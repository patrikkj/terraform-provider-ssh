package provider

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/sshconf"
)

var _ resource.Resource = &SSHConfigResource{}

func NewSSHConfigResource() resource.Resource {
	return &SSHConfigResource{}
}

type SSHConfigResource struct{}

func (r *SSHConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *SSHConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SSHConfigResourceSchema
}

func (r *SSHConfigResource) Configure(_ context.Context, _ resource.ConfigureRequest, _ *resource.ConfigureResponse) {
}

func (r *SSHConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHConfigResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique ID based on the path and timestamp
	data.Id = types.StringValue(generateFileID(data.Path.ValueString(), time.Now()))

	// Expand the path (handle ~)
	path := data.Path.ValueString()
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			resp.Diagnostics.AddError("Failed to expand home directory", err.Error())
			return
		}
		path = filepath.Join(home, path[1:])
	}

	// Read or create the config file
	config, err := sshconf.ParseConfigFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			resp.Diagnostics.AddError("Failed to read SSH config file", err.Error())
			return
		}
		// Create new empty config if file doesn't exist
		config = sshconf.ParseConfig("")
	}

	// Apply the patch
	err = config.Patch(data.Find.ValueString(), data.Patch.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to apply SSH config patch", err.Error())
		return
	}

	// Write the config back to file
	err = config.WriteFile(path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to write SSH config file", err.Error())
		return
	}

	// Update the content and lines in the model
	content := config.Render()

	// Convert the lines to our model format
	lines := make([]SSHConfigLine, len(config.Lines()))
	for i, line := range config.Lines() {
		lines[i] = convertSSHConfLine(line)
	}

	lineValues := make([]attr.Value, len(lines))
	for i, line := range lines {
		lineValues[i] = line.toAttrValue()
	}

	linesList, diags := types.ListValue(
		sshConfigLineObjectType,
		lineValues,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Content = types.StringValue(content)
	data.Lines = linesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSHConfigResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Expand the path (handle ~)
	path := data.Path.ValueString()
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			resp.Diagnostics.AddError("Failed to expand home directory", err.Error())
			return
		}
		path = filepath.Join(home, path[1:])
	}

	// Read the config file
	config, err := sshconf.ParseConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read SSH config file", err.Error())
		return
	}

	// Update the content and lines in the model
	content := config.Render()

	// Convert the lines to our model format
	lines := make([]SSHConfigLine, len(config.Lines()))
	for i, line := range config.Lines() {
		lines[i] = convertSSHConfLine(line)
	}

	lineValues := make([]attr.Value, len(lines))
	for i, line := range lines {
		lineValues[i] = line.toAttrValue()
	}

	linesList, diags := types.ListValue(
		sshConfigLineObjectType,
		lineValues,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Content = types.StringValue(content)
	data.Lines = linesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSHConfigResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Expand the path (handle ~)
	path := data.Path.ValueString()
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			resp.Diagnostics.AddError("Failed to expand home directory", err.Error())
			return
		}
		path = filepath.Join(home, path[1:])
	}

	// Read the config file
	config, err := sshconf.ParseConfigFile(path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read SSH config file", err.Error())
		return
	}

	// Apply the patch
	err = config.Patch(data.Find.ValueString(), data.Patch.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to apply SSH config patch", err.Error())
		return
	}

	// Write the config back to file
	err = config.WriteFile(path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to write SSH config file", err.Error())
		return
	}

	// Update the content and lines in the model
	content := config.Render()

	// Convert the lines to our model format
	lines := make([]SSHConfigLine, len(config.Lines()))
	for i, line := range config.Lines() {
		lines[i] = convertSSHConfLine(line)
	}

	lineValues := make([]attr.Value, len(lines))
	for i, line := range lines {
		lineValues[i] = line.toAttrValue()
	}

	linesList, diags := types.ListValue(
		sshConfigLineObjectType,
		lineValues,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Content = types.StringValue(content)
	data.Lines = linesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SSHConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSHConfigResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If delete_on_destroy is false, we're done
	if !data.DeleteOnDestroy.ValueBool() {
		return
	}

	// Expand the path (handle ~)
	path := data.Path.ValueString()
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			resp.Diagnostics.AddError("Failed to expand home directory", err.Error())
			return
		}
		path = filepath.Join(home, path[1:])
	}

	// Read the config file
	config, err := sshconf.ParseConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		resp.Diagnostics.AddError("Failed to read SSH config file", err.Error())
		return
	}

	// Delete the patched section
	err = config.Delete(data.Find.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete SSH config section", err.Error())
		return
	}

	// Write the config back to file
	err = config.WriteFile(path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to write SSH config file", err.Error())
		return
	}
}
