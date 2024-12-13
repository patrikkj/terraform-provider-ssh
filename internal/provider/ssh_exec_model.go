package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
