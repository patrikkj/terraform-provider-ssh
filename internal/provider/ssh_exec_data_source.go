package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = SSHExecDataSourceSchema
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
	var data SSHExecDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default value for fail_if_nonzero if not specified
	if data.FailIfNonzero.IsNull() {
		data.FailIfNonzero = types.BoolValue(true)
	}

	// Generate ID early, based on the command
	data.Id = types.StringValue(generateExecID(data.Command.ValueString(), time.Now()))

	// Get SSH client
	client, newClient, err := d.manager.GetClient(
		*data.SSHConnectionModel.toConfig(),
		data.UseProviderAsBastion.ValueBool(),
		data.Bastion.toConfig(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
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
