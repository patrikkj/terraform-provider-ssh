package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
