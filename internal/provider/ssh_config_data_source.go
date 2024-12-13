package provider

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/sshconf"
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

	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read SSH config file", err.Error())
		return
	}

	// Parse the SSH config
	config := sshconf.ParseConfig(string(content))

	// Convert the lines to our model format
	lines := make([]SSHConfigLine, len(config.Lines()))
	for i, line := range config.Lines() {
		lines[i] = convertSSHConfLine(line)
	}

	// Convert the lines to attr.Value
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

	// Set the values
	data.Content = types.StringValue(string(content))
	data.Lines = linesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
