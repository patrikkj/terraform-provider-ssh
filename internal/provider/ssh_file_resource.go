package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/sftp"
)

var _ resource.Resource = &SSHFileResource{}

func NewSSHFileResource() resource.Resource {
	return &SSHFileResource{}
}

type SSHFileResource struct {
	manager *SSHManager
}

type SSHFileResourceModel struct {
	Path                 types.String `tfsdk:"path"`
	Content              types.String `tfsdk:"content"`
	Permissions          types.String `tfsdk:"permissions"`
	Id                   types.String `tfsdk:"id"`
	Host                 types.String `tfsdk:"host"`
	User                 types.String `tfsdk:"user"`
	Password             types.String `tfsdk:"password"`
	PrivateKey           types.String `tfsdk:"private_key"`
	UseProviderAsBastion types.Bool   `tfsdk:"use_provider_as_bastion"`
}

func (r *SSHFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *SSHFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Create specific attributes
	attributes := map[string]schema.Attribute{
		"path": schema.StringAttribute{
			MarkdownDescription: "Path to the file",
			Required:            true,
		},
		"content": schema.StringAttribute{
			MarkdownDescription: "Content to write to the file",
			Required:            true,
		},
		"permissions": schema.StringAttribute{
			MarkdownDescription: "File permissions (e.g., '0644')",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("0644"),
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this file",
			Computed:            true,
		},
	}

	// Merge with common SSH connection attributes
	for k, v := range GetCommonSSHConnectionSchema() {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage files over SSH",
		Attributes:          attributes,
	}
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

func generateFileID(path string) string {
	h := md5.New()
	h.Write([]byte(path))
	return hex.EncodeToString(h.Sum(nil))
}

func (r *SSHFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSHFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique, stable ID using the file path
	data.Id = types.StringValue(generateFileID(data.Path.ValueString()))

	if err := r.writeFile(ctx, &data); err != nil {
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

	// Create SFTP client
	sftpClient, err := sftp.NewClient(r.manager.client)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SFTP client", err.Error())
		return
	}
	defer sftpClient.Close()

	// Open and read the file
	f, err := sftpClient.Open(data.Path.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read file", err.Error())
		return
	}

	// Strip trailing newline when reading
	contentStr := string(content)
	if len(contentStr) > 0 && strings.HasSuffix(contentStr, "\n") {
		contentStr = contentStr[:len(contentStr)-1]
	}
	data.Content = types.StringValue(contentStr)

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

	// Write the file
	if err := r.writeFile(ctx, &data); err != nil {
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

	// Create SFTP client
	sftpClient, err := sftp.NewClient(r.manager.client)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SFTP client", err.Error())
		return
	}
	defer sftpClient.Close()

	if err := sftpClient.Remove(data.Path.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to delete file", err.Error())
		return
	}
}

func (r *SSHFileResource) writeFile(ctx context.Context, data *SSHFileResourceModel) error {
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

	// Create SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Create directory if needed
	dirPath := filepath.Dir(data.Path.ValueString())
	if err := sftpClient.MkdirAll(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or truncate the file
	f, err := sftpClient.Create(data.Path.ValueString())
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Write content
	content := data.Content.ValueString()
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if _, err := f.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	// Close the file before changing permissions
	f.Close()

	// Try to set permissions, but don't fail if it doesn't work
	mode := parseFileMode(data.Permissions.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Attempting to chmod %s to %s", data.Path.ValueString(), mode))

	if err := sftpClient.Chmod(data.Path.ValueString(), mode); err != nil {
		// Log the permission error but don't fail the resource creation/update
		tflog.Warn(ctx, fmt.Sprintf("Warning: Could not set permissions on %s to %s: %s",
			data.Path.ValueString(), mode, err))
	}

	// Only generate a new ID if one hasn't been set
	if data.Id.IsNull() {
		data.Id = types.StringValue(generateFileID(data.Path.ValueString()))
	}

	return nil
}

// Helper function to parse file mode
func parseFileMode(mode string) fs.FileMode {
	var result uint32
	if _, err := fmt.Sscanf(mode, "%o", &result); err != nil {
		return 0644 // Default to 0644 if parsing fails
	}
	return fs.FileMode(result)
}
