package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SSHConfigDataSource{}

func NewSSHConfigDataSource() datasource.DataSource {
	return &SSHConfigDataSource{}
}

type SSHConfigDataSource struct{}

func (d *SSHConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *SSHConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SSHConfigDataSourceSchema
}

func (d *SSHConfigDataSource) Configure(_ context.Context, _ datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
}

func (d *SSHConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHConfigDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a unique ID based on the path and timestamp
	data.Id = types.StringValue(generateFileID(data.Path.ValueString(), time.Now()))

	config, err := readOrCreateConfig(data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read SSH config", err.Error())
		return
	}

	// Update model content and lines
	content, lines, err := updateModelFromConfig(config)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update model", err.Error())
		return
	}

	data.Content = content
	data.Lines = lines

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
