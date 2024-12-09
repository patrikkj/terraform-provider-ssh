package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = SSHFileDataSourceSchema
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

	// Generate a unique ID early, based on the path
	data.Id = types.StringValue(generateFileID(data.Path.ValueString(), time.Now()))

	client, newClient, err := d.manager.GetClient(&data.SSHConnectionModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SSH client", err.Error())
		return
	}

	if newClient {
		defer client.Close()
	}

	content, err := readFile(client, data.Path.ValueString())
	if err != nil {
		if data.FailIfAbsent.ValueBool() {
			resp.Diagnostics.AddError("Failed to read file", err.Error())
			return
		}
		// If fail_if_absent is false, return empty content
		data.Content = types.StringValue("")
	} else {
		data.Content = types.StringValue(content)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
