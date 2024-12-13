package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHConfigLine struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Indent      types.String `tfsdk:"indent"`
	Sep         types.String `tfsdk:"sep"`
	Comment     types.String `tfsdk:"comment"`
	TrailIndent types.String `tfsdk:"trail_indent"`
	Children    types.List   `tfsdk:"children"`
}

type SSHConfigDataSourceModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
	Lines   types.List   `tfsdk:"lines"`
	Id      types.String `tfsdk:"id"`
}

type SSHConfigResourceModel struct {
	Path            types.String `tfsdk:"path"`
	Content         types.String `tfsdk:"content"`
	Patch           types.String `tfsdk:"patch"`
	Find            types.String `tfsdk:"find"`
	DeleteOnDestroy types.Bool   `tfsdk:"delete_on_destroy"`
	Lines           types.List   `tfsdk:"lines"`
	Id              types.String `tfsdk:"id"`
}
