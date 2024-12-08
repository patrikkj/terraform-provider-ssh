package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

var _ resource.Resource = &SSHExecResource{}

func NewSSHExecResource() resource.Resource {
	return &SSHExecResource{}
}

type SSHExecResource struct {
	client *ssh.Client
}

type SSHExecResourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	Id            types.String `tfsdk:"id"`
}

func (r *SSHExecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exec"
}

func (r *SSHExecResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Execute commands over SSH with potential side effects",
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				MarkdownDescription: "Command to execute",
				Required:            true,
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "Output of the command",
				Computed:            true,
			},
			"exit_code": schema.Int64Attribute{
				MarkdownDescription: "Exit code of the command",
				Computed:            true,
			},
			"fail_if_nonzero": schema.BoolAttribute{
				MarkdownDescription: "Whether to fail if the command returns a non-zero exit code. Defaults to true if not specified.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this execution",
				Computed:            true,
			},
		},
	}
}

func (r *SSHExecResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ssh.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ssh.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SSHExecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHExecResourceModel

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
	if data.Id.IsNull() {
		data.Id = types.StringValue(data.Command.ValueString())
	}

	if err := r.executeCommand(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHExecResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No need to re-run the command during read
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHExecResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.executeCommand(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No action needed for delete
}

func (r *SSHExecResource) executeCommand(ctx context.Context, data *SSHExecResourceModel) error {
	// Set default values for computed fields
	if data.Output.IsNull() {
		data.Output = types.StringValue("")
	}
	if data.ExitCode.IsNull() {
		data.ExitCode = types.Int64Value(0)
	}
	if data.Id.IsNull() {
		data.Id = types.StringValue(data.Command.ValueString())
	}

	session, err := r.client.NewSession()
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
	data.Id = types.StringValue(data.Command.ValueString())

	return nil
}
