package provider

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var _ datasource.DataSource = &SSHFileDataSource{}

func NewSSHFileDataSource() datasource.DataSource {
	return &SSHFileDataSource{}
}

type SSHFileDataSource struct {
	client *ssh.Client
}

type SSHFileDataSourceModel struct {
	Path         types.String `tfsdk:"path"`
	Content      types.String `tfsdk:"content"`
	FailIfAbsent types.Bool   `tfsdk:"fail_if_absent"`
	Id           types.String `tfsdk:"id"`
}

func (d *SSHFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *SSHFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read files over SSH",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Path to the file to read",
				Required:            true,
			},
			"fail_if_absent": schema.BoolAttribute{
				MarkdownDescription: "Fail if the file does not exist",
				Optional:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Content of the file",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this file read",
				Computed:            true,
			},
		},
	}
}

func (d *SSHFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ssh.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ssh.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *SSHFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHFileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(d.client)
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
