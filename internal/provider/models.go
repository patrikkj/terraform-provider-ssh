package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Common model for SSH connection configuration
type SSHConnectionModel struct {
	Host                 types.String `tfsdk:"host"`
	User                 types.String `tfsdk:"user"`
	Password             types.String `tfsdk:"password"`
	PrivateKey           types.String `tfsdk:"private_key"`
	UseProviderAsBastion types.Bool   `tfsdk:"use_provider_as_bastion"`
	BastionHost          types.String `tfsdk:"bastion_host"`
	BastionPort          types.Int64  `tfsdk:"bastion_port"`
	BastionUser          types.String `tfsdk:"bastion_user"`
	BastionPassword      types.String `tfsdk:"bastion_password"`
	BastionPrivateKey    types.String `tfsdk:"bastion_private_key"`
}

type SSHProviderModel struct {
	Host              types.String `tfsdk:"host"`
	Port              types.Int64  `tfsdk:"port"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	PrivateKey        types.String `tfsdk:"private_key"`
	BastionHost       types.String `tfsdk:"bastion_host"`
	BastionPort       types.Int64  `tfsdk:"bastion_port"`
	BastionUser       types.String `tfsdk:"bastion_user"`
	BastionPassword   types.String `tfsdk:"bastion_password"`
	BastionPrivateKey types.String `tfsdk:"bastion_private_key"`
}

type SSHExecDataSourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	Id            types.String `tfsdk:"id"`
	SSHConnectionModel
}

type SSHExecResourceModel struct {
	Command       types.String `tfsdk:"command"`
	Output        types.String `tfsdk:"output"`
	ExitCode      types.Int64  `tfsdk:"exit_code"`
	FailIfNonzero types.Bool   `tfsdk:"fail_if_nonzero"`
	OnDestroy     types.String `tfsdk:"on_destroy"`
	Id            types.String `tfsdk:"id"`
	SSHConnectionModel
}

type SSHFileDataSourceModel struct {
	Path         types.String `tfsdk:"path"`
	Content      types.String `tfsdk:"content"`
	Permissions  types.String `tfsdk:"permissions"`
	FailIfAbsent types.Bool   `tfsdk:"fail_if_absent"`
	Id           types.String `tfsdk:"id"`
	SSHConnectionModel
}

type SSHFileResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Content         types.String `tfsdk:"content"`
	Permissions     types.String `tfsdk:"permissions"`
	FailIfAbsent    types.Bool   `tfsdk:"fail_if_absent"`
	DeleteOnDestroy types.Bool   `tfsdk:"delete_on_destroy"`
	Id              types.String `tfsdk:"id"`
	SSHConnectionModel
}
