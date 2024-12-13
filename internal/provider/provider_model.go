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
	// Handle nil receiver
	if m == nil {
		return nil
	}

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
