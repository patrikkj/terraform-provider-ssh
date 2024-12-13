package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

var SSHConfigLineAttrs = struct {
	Key         dschema.StringAttribute
	Value       dschema.StringAttribute
	Indent      dschema.StringAttribute
	Sep         dschema.StringAttribute
	Comment     dschema.StringAttribute
	TrailIndent dschema.StringAttribute
	Children    dschema.ListNestedAttribute
}{
	Key: dschema.StringAttribute{
		MarkdownDescription: "The key/directive of the config line",
		Computed:            true,
	},
	Value: dschema.StringAttribute{
		MarkdownDescription: "The value of the config line",
		Computed:            true,
	},
	Indent: dschema.StringAttribute{
		MarkdownDescription: "The indentation before the key",
		Computed:            true,
	},
	Sep: dschema.StringAttribute{
		MarkdownDescription: "The separator between key and value",
		Computed:            true,
	},
	Comment: dschema.StringAttribute{
		MarkdownDescription: "Any comment on the line",
		Computed:            true,
	},
	TrailIndent: dschema.StringAttribute{
		MarkdownDescription: "Any trailing indentation",
		Computed:            true,
	},
}

var SSHConfigChildrenAttr = dschema.ListNestedAttribute{
	MarkdownDescription: "SSH config children",
	Optional:            true,
	NestedObject: dschema.NestedAttributeObject{
		Attributes: map[string]dschema.Attribute{
			"key":          SSHConfigLineAttrs.Key,
			"value":        SSHConfigLineAttrs.Value,
			"indent":       SSHConfigLineAttrs.Indent,
			"sep":          SSHConfigLineAttrs.Sep,
			"comment":      SSHConfigLineAttrs.Comment,
			"trail_indent": SSHConfigLineAttrs.TrailIndent,
			// Add dummy children, make it an empty attribute
			"children": dschema.ListNestedAttribute{
				MarkdownDescription: "SSH config children",
				Optional:            true,
				NestedObject:        dschema.NestedAttributeObject{},
			},
		},
	},
}

var SSHConfigLineSchema = dschema.NestedAttributeObject{
	Attributes: map[string]dschema.Attribute{
		"key":          SSHConfigLineAttrs.Key,
		"value":        SSHConfigLineAttrs.Value,
		"indent":       SSHConfigLineAttrs.Indent,
		"sep":          SSHConfigLineAttrs.Sep,
		"comment":      SSHConfigLineAttrs.Comment,
		"trail_indent": SSHConfigLineAttrs.TrailIndent,
		"children":     SSHConfigChildrenAttr,
	},
}

var SSHConfigAttrs = struct {
	Path            dschema.StringAttribute
	Content         dschema.StringAttribute
	Lines           dschema.ListNestedAttribute
	ID              dschema.StringAttribute
	Patch           dschema.StringAttribute
	Find            dschema.StringAttribute
	DeleteOnDestroy rschema.BoolAttribute
}{
	Path: dschema.StringAttribute{
		MarkdownDescription: "Path to the SSH config file",
		Required:            true,
	},
	Content: dschema.StringAttribute{
		MarkdownDescription: "Raw content of the SSH config file",
		Computed:            true,
	},
	Lines: dschema.ListNestedAttribute{
		MarkdownDescription: "Parsed SSH config lines",
		Computed:            true,
		NestedObject:        SSHConfigLineSchema,
	},
	ID: dschema.StringAttribute{
		MarkdownDescription: "Unique identifier for this SSH config",
		Computed:            true,
	},
	Patch: dschema.StringAttribute{
		MarkdownDescription: "The SSH config patch to apply",
		Required:            true,
	},
	Find: dschema.StringAttribute{
		MarkdownDescription: "The line to find and patch (e.g., 'Host example')",
		Required:            true,
	},
	DeleteOnDestroy: rschema.BoolAttribute{
		MarkdownDescription: "Whether to delete the patched section when the resource is destroyed",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
}

var SSHConfigDataSourceSchema = dschema.Schema{
	MarkdownDescription: "Read and parse SSH config files",
	Attributes: map[string]dschema.Attribute{
		"path":    SSHConfigAttrs.Path,
		"content": SSHConfigAttrs.Content,
		"lines":   SSHConfigAttrs.Lines,
		"id":      SSHConfigAttrs.ID,
	},
}

var SSHConfigResourceSchema = rschema.Schema{
	MarkdownDescription: "Manage SSH config files",
	Attributes: map[string]rschema.Attribute{
		"path":              SSHConfigAttrs.Path,
		"content":           SSHConfigAttrs.Content,
		"patch":             SSHConfigAttrs.Patch,
		"find":              SSHConfigAttrs.Find,
		"delete_on_destroy": SSHConfigAttrs.DeleteOnDestroy,
		"lines":             SSHConfigAttrs.Lines,
		"id":                SSHConfigAttrs.ID,
	},
}
