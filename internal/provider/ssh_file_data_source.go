package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/sftp"
)

var _ datasource.DataSource = &SSHFileDataSource{}

func NewSSHFileDataSource() datasource.DataSource {
	return &SSHFileDataSource{}
}

type SSHFileDataSource struct {
	manager *SSHManager
}

func (d *SSHFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *SSHFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// Create specific attributes
	attributes := map[string]schema.Attribute{
		"path": schema.StringAttribute{
			MarkdownDescription: "Path to the file to read",
			Required:            true,
		},
		"content": schema.StringAttribute{
			MarkdownDescription: "Content of the file",
			Computed:            true,
		},
		"fail_if_absent": schema.BoolAttribute{
			MarkdownDescription: "Fail if the file does not exist",
			Optional:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this file read",
			Computed:            true,
		},
	}

	// Merge with common SSH connection attributes
	for k, v := range GetCommonSSHConnectionSchema() {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Read files over SSH",
		Attributes:          attributes,
	}
}

func (d *SSHFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	manager, ok := req.ProviderData.(*SSHManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *SSHManager, got: %T", req.ProviderData),
		)
		return
	}

	d.manager = manager
}

func (d *SSHFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHFileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, newClient, err := d.manager.GetClient(&SSHConnectionConfig{
		Host:                 data.Host,
		User:                 data.User,
		Password:             data.Password,
		PrivateKey:           data.PrivateKey,
		UseProviderAsBastion: data.UseProviderAsBastion,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	// If we created a new client, close it when done
	if newClient {
		defer client.Close()
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SFTP client", err.Error())
		return
	}
	defer sftpClient.Close()

	// Open and read the file
	f, err := sftpClient.Open(data.Path.ValueString())
	if err != nil {
		if data.FailIfAbsent.ValueBool() {
			resp.Diagnostics.AddError("Failed to read file", err.Error())
			return
		}
		// If fail_if_absent is false, return empty content
		data.Content = types.StringValue("")
	} else {
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			resp.Diagnostics.AddError("Failed to read file contents", err.Error())
			return
		}
		data.Content = types.StringValue(string(content))
	}

	data.Id = types.StringValue(data.Path.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
