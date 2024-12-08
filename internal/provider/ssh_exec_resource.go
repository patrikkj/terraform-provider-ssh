package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

var _ resource.Resource = &SSHExecResource{}

func NewSSHExecResource() resource.Resource {
	return &SSHExecResource{}
}

type SSHExecResource struct {
	manager *SSHManager
}

func (r *SSHExecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exec"
}

func (r *SSHExecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Execute commands over SSH with potential side effects",
		Attributes:          GetSSHExecSchema(),
	}
}

func (r *SSHExecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	manager, ok := req.ProviderData.(*SSHManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *SSHManager, got: %T", req.ProviderData),
		)
		return
	}

	r.manager = manager
}

func generateExecID(command string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(command))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}

func (r *SSHExecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHExecModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default values for computed fields
	if data.Output.IsNull() {
		data.Output = types.StringValue("")
	}
	if data.ExitCode.IsNull() {
		data.ExitCode = types.Int64Value(0)
	}

	// Generate a unique, stable ID using command and creation timestamp
	timestamp := time.Now()
	data.Id = types.StringValue(generateExecID(data.Command.ValueString(), timestamp))

	if err := r.executeCommand(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHExecModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No need to re-run the command during read
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHExecModel

	// Get the current state
	var state SSHExecModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the planned changes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the original ID from state
	data.Id = state.Id

	// Execute the command
	if err := r.executeCommand(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHExecModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there's an on_destroy command, execute it
	if !data.OnDestroy.IsNull() {
		// Create temporary model for the destroy command
		destroyData := SSHExecModel{
			Command:       data.OnDestroy,
			FailIfNonzero: data.FailIfNonzero,
		}

		if err := r.executeCommand(ctx, &destroyData); err != nil {
			resp.Diagnostics.AddError("Failed to execute destroy command", err.Error())
			return
		}
	}
}

func (r *SSHExecResource) executeCommand(ctx context.Context, data *SSHExecModel) error {
	client, newClient, err := r.manager.GetClient(&SSHConnectionConfig{
		Host:                 data.Host,
		User:                 data.User,
		Password:             data.Password,
		PrivateKey:           data.PrivateKey,
		UseProviderAsBastion: data.UseProviderAsBastion,
	})
	if err != nil {
		return fmt.Errorf("failed to get SSH client: %w", err)
	}

	// If we created a new client, close it when done
	if newClient {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Execute command and capture output
	output, err := session.CombinedOutput(data.Command.ValueString())
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode := int64(exitErr.ExitStatus())
			data.ExitCode = types.Int64Value(exitCode)
			if data.FailIfNonzero.ValueBool() && exitCode != 0 {
				return fmt.Errorf("command exited with non-zero status: %d", exitCode)
			}
		} else {
			return err
		}
	} else {
		data.ExitCode = types.Int64Value(0)
	}

	data.Output = types.StringValue(string(output))

	return nil
}
