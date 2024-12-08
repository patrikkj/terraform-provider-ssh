package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

var _ datasource.DataSource = &SSHExecDataSource{}

func NewSSHExecDataSource() datasource.DataSource {
	return &SSHExecDataSource{}
}

type SSHExecDataSource struct {
	client *ssh.Client
}

type SSHExecDataSourceModel struct {
	Command       types.String `tfsdk:"command"`
	Stdout        types.String `tfsdk:"stdout"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	Id            types.String `tfsdk:"id"`
}

func (d *SSHExecDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exec"
}

func (d *SSHExecDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Execute commands over SSH",
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				MarkdownDescription: "Command to execute",
				Required:            true,
			},
			"stdout": schema.StringAttribute{
				MarkdownDescription: "Output of the command",
				Computed:            true,
			},
			"exit_code": schema.Int64Attribute{
				MarkdownDescription: "Exit code of the command",
				Computed:            true,
			},
			"fail_if_nonzero": schema.BoolAttribute{
				MarkdownDescription: "Whether to fail if the command returns a non-zero exit code",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this execution",
				Computed:            true,
			},
		},
	}
}

func (d *SSHExecDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SSHExecDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHExecDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default value for fail_if_nonzero if not specified
	if data.FailIfNonzero.IsNull() {
		data.FailIfNonzero = types.BoolValue(true)
	}

	session, err := d.client.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SSH session", err.Error())
		return
	}
	defer session.Close()

	// Execute command and capture output
	output, err := session.CombinedOutput(data.Command.ValueString())
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode := int64(exitErr.ExitStatus())
			data.ExitCode = types.Int64Value(exitCode)
			data.Stdout = types.StringValue(string(output))
			if data.FailIfNonzero.ValueBool() && exitCode != 0 {
				resp.Diagnostics.AddError("Command exited with non-zero status",
					fmt.Sprintf("Exit code: %d, Output: %s", exitCode, output))
				return
			}
		} else {
			resp.Diagnostics.AddError("Failed to execute command", err.Error())
			return
		}
	} else {
		data.ExitCode = types.Int64Value(0)
		data.Stdout = types.StringValue(string(output))
	}

	data.Id = types.StringValue(fmt.Sprintf("%s-%d", data.Command.ValueString(), data.ExitCode.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
