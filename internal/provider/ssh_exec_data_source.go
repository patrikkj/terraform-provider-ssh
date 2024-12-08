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
	manager *SSHManager
}

func (d *SSHExecDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exec"
}

func (d *SSHExecDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// Create specific attributes
	attributes := map[string]schema.Attribute{
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
			MarkdownDescription: "Whether to fail if the command returns a non-zero exit code",
			Optional:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this execution",
			Computed:            true,
		},
	}

	// Merge with common SSH connection attributes
	for k, v := range GetCommonSSHConnectionSchema() {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Execute commands over SSH",
		Attributes:          attributes,
	}
}

func (d *SSHExecDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SSHExecDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHExecModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default value for fail_if_nonzero if not specified
	if data.FailIfNonzero.IsNull() {
		data.FailIfNonzero = types.BoolValue(true)
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

	session, err := client.NewSession()
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
			data.Output = types.StringValue(string(output))
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
		data.Output = types.StringValue(string(output))
	}

	data.Id = types.StringValue(fmt.Sprintf("%s-%d", data.Command.ValueString(), data.ExitCode.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
