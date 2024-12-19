package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SSHConfigLineChild represents a child line without nested children
type SSHConfigLineChild struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Indent      types.String `tfsdk:"indent"`
	Sep         types.String `tfsdk:"sep"`
	Comment     types.String `tfsdk:"comment"`
	TrailIndent types.String `tfsdk:"trail_indent"`
}

// SSHConfigLine represents a line that may have children
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

type dSSHConfigResourceValues struct {
	Id              string
	Path            string
	Find            string
	Patch           string
	Content         string
	DeleteOnDestroy bool
	Lines           []SSHConfigLine
}

func (m *SSHConfigResourceModel) toValues() dSSHConfigResourceValues {
	var lines []SSHConfigLine
	if !m.Lines.IsNull() {
		m.Lines.ElementsAs(context.Background(), &lines, false)
	}

	// Use the default value (true) if DeleteOnDestroy is null
	deleteOnDestroy := true
	if !m.DeleteOnDestroy.IsNull() {
		deleteOnDestroy = m.DeleteOnDestroy.ValueBool()
	}

	return dSSHConfigResourceValues{
		Id:              m.Id.ValueString(),
		Path:            m.Path.ValueString(),
		Find:            m.Find.ValueString(),
		Patch:           m.Patch.ValueString(),
		Content:         m.Content.ValueString(),
		DeleteOnDestroy: deleteOnDestroy,
		Lines:           lines,
	}
}
