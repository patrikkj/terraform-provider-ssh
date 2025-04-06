package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHProviderModel struct {
	SSHConnectionModel
	Bastion *SSHConnectionModel `tfsdk:"bastion"`
}

var SSHProviderSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"host":        schema.StringAttribute{Description: "The hostname or IP address of the target SSH server", Optional: true},
		"port":        schema.Int64Attribute{Description: "The port number of the target SSH server", Optional: true},
		"user":        schema.StringAttribute{Description: "The username for SSH authentication", Optional: true},
		"password":    schema.StringAttribute{Description: "The password for SSH authentication", Optional: true, Sensitive: true},
		"private_key": schema.StringAttribute{Description: "The private key for SSH authentication", Optional: true, Sensitive: true},
		"bastion": schema.SingleNestedAttribute{
			Description: "Bastion host configuration",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"host":        SSHConnectionSchema.Host,
				"port":        SSHConnectionSchema.Port,
				"user":        SSHConnectionSchema.User,
				"password":    SSHConnectionSchema.Password,
				"private_key": SSHConnectionSchema.PrivateKey,
			},
		},
	},
}

var _ provider.Provider = &SSHProvider{}

type SSHProvider struct {
	version string
	manager *SSHManager
}

func (p *SSHProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ssh"
	resp.Version = p.version
}

func (p *SSHProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = SSHProviderSchema
}

func (p *SSHProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SSHProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default port if not specified
	if config.Port.IsNull() {
		config.Port = types.Int64Value(22)
	}

	// Create the SSH manager with provider configuration
	manager, err := NewSSHManager(config.SSHConnectionModel.toConfig(), config.Bastion.toConfig())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create SSH manager",
			err.Error(),
		)
		return
	}

	p.manager = manager
	resp.DataSourceData = p.manager
	resp.ResourceData = p.manager
}

func (p *SSHProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHExecResource,
		NewSSHFileResource,
	}
}

func (p *SSHProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSSHExecDataSource,
		NewSSHFileDataSource,
	}
}

func (p *SSHProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SSHProvider{
			version: version,
		}
	}
}
