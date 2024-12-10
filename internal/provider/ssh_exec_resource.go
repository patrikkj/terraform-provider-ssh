package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = SSHExecResourceSchema
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

	// Generate a unique, stable ID before executing the command
	data.Id = types.StringValue(generateExecID(data.Command.ValueString(), time.Now()))

	// Get SSH client
	client, err := r.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	// Execute the command
	output, exitCode, err := executeCommand(
		client,
		data.Command.ValueString(),
		data.FailIfNonzero.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	data.Output = types.StringValue(output)
	data.ExitCode = types.Int64Value(exitCode)

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

	// Get the current state
	var state SSHExecResourceModel
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

	// Get SSH client
	client, err := r.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	// Execute the command
	output, exitCode, err := executeCommand(
		client,
		data.Command.ValueString(),
		data.FailIfNonzero.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}
	data.Output = types.StringValue(output)
	data.ExitCode = types.Int64Value(exitCode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHExecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHExecResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there's an on_destroy command, execute it
	if !data.OnDestroy.IsNull() {
		// Get SSH client
		client, err := r.manager.GetClient(
			*data.SSHConnectionModel.toConfig(),
			data.UseProviderAsBastion.ValueBool(),
			data.Bastion.toConfig(),
			nil,
		)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
			return
		}

		_, _, err = executeCommand(
			client,
			data.OnDestroy.ValueString(),
			data.FailIfNonzero.ValueBool(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Failed to execute destroy command", err.Error())
			return
		}
	}
}
