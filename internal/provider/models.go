package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Common model for SSH connection configuration
type SSHConnectionModel struct {
	Host       types.String `tfsdk:"host"`
	User       types.String `tfsdk:"user"`
	Password   types.String `tfsdk:"password"`
	PrivateKey types.String `tfsdk:"private_key"`
	Port       types.Int64  `tfsdk:"port"`
}

type SSHConnectionConfig struct {
	Host       *string
	User       *string
	Password   *string
	PrivateKey *string
	Port       *int64
}

func (m *SSHConnectionModel) toConfig() *SSHConnectionConfig {
	config := &SSHConnectionConfig{}

	if !m.Host.IsNull() {
		value := m.Host.ValueString()
		config.Host = &value
	}
	if !m.User.IsNull() {
		value := m.User.ValueString()
		config.User = &value
	}
	if !m.Password.IsNull() {
		value := m.Password.ValueString()
		config.Password = &value
	}
	if !m.PrivateKey.IsNull() {
		value := m.PrivateKey.ValueString()
		config.PrivateKey = &value
	}
	if !m.Port.IsNull() {
		value := m.Port.ValueInt64()
		config.Port = &value
	}

	return config
}

type SSHProviderModel struct {
	SSHConnectionModel
	Bastion *SSHConnectionModel `tfsdk:"bastion"`
}

type SSHExecDataSourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	Id            types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

type SSHExecResourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	OnDestroy     types.String `tfsdk:"on_destroy"`
	Id            types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

type SSHFileDataSourceModel struct {
	Path         types.String `tfsdk:"path"`
	Content      types.String `tfsdk:"content"`
	Permissions  types.String `tfsdk:"permissions"`
	FailIfAbsent types.Bool   `tfsdk:"fail_if_absent"`
	Id           types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}

type SSHFileResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Content         types.String `tfsdk:"content"`
	Permissions     types.String `tfsdk:"permissions"`
	FailIfAbsent    types.Bool   `tfsdk:"fail_if_absent"`
	DeleteOnDestroy types.Bool   `tfsdk:"delete_on_destroy"`
	Id              types.String `tfsdk:"id"`

	// Connection details
	SSHConnectionModel
	UseProviderAsBastion types.Bool          `tfsdk:"use_provider_as_bastion"`
	Bastion              *SSHConnectionModel `tfsdk:"bastion"`
}
