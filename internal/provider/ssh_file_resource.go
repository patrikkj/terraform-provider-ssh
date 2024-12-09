package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &SSHFileResource{}

func NewSSHFileResource() resource.Resource {
	return &SSHFileResource{}
}

type SSHFileResource struct {
	manager *SSHManager
}

func (r *SSHFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *SSHFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SSHFileResourceSchema
}

func (r *SSHFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SSHFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique, stable ID using the file path
	data.Id = types.StringValue(generateFileID(data.Path.ValueString()))

	client, newClient, err := r.manager.GetClient(&data.SSHConnectionModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
	}

	if err := writeFile(ctx, client, data.Path.ValueString(), data.Content.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to write file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, newClient, err := r.manager.GetClient(&data.SSHConnectionModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
	}

	content, err := readFile(client, data.Path.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Content = types.StringValue(content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHFileResourceModel

	// Get the current state
	var state SSHFileResourceModel
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

	client, newClient, err := r.manager.GetClient(&data.SSHConnectionModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
	}

	if err := writeFile(ctx, client, data.Path.ValueString(), data.Content.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to update file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if delete_on_destroy is set to false
	if !data.DeleteOnDestroy.IsNull() && !data.DeleteOnDestroy.ValueBool() {
		// Skip deletion if delete_on_destroy is false
		return
	}

	client, newClient, err := r.manager.GetClient(&data.SSHConnectionModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
	}

	if err := deleteFile(client, data.Path.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete file", err.Error())
		return
	}
}
