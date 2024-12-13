package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

var SSHConfigAttrs = struct {
	Key         dschema.StringAttribute
	Value       dschema.StringAttribute
	Indent      dschema.StringAttribute
	Sep         dschema.StringAttribute
	Comment     dschema.StringAttribute
	TrailIndent dschema.StringAttribute
	Children    dschema.ListNestedAttribute

	Path            dschema.StringAttribute
	Content         dschema.StringAttribute
	ID              dschema.StringAttribute
	Patch           dschema.StringAttribute
	Find            dschema.StringAttribute
	DeleteOnDestroy rschema.BoolAttribute
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
	Path: dschema.StringAttribute{
		MarkdownDescription: "Path to the SSH config file",
		Required:            true,
	},
	Content: dschema.StringAttribute{
		MarkdownDescription: "Raw content of the SSH config file",
		Computed:            true,
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

var SSHConfigLineSchema = dschema.ListNestedAttribute{
	MarkdownDescription: "Parsed SSH config lines",
	Computed:            true,
	NestedObject: dschema.NestedAttributeObject{
		Attributes: map[string]dschema.Attribute{
			"key":          SSHConfigAttrs.Key,
			"value":        SSHConfigAttrs.Value,
			"indent":       SSHConfigAttrs.Indent,
			"sep":          SSHConfigAttrs.Sep,
			"comment":      SSHConfigAttrs.Comment,
			"trail_indent": SSHConfigAttrs.TrailIndent,
			"children": dschema.ListNestedAttribute{
				MarkdownDescription: "SSH config children",
				Optional:            true,
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"key":          SSHConfigAttrs.Key,
						"value":        SSHConfigAttrs.Value,
						"indent":       SSHConfigAttrs.Indent,
						"sep":          SSHConfigAttrs.Sep,
						"comment":      SSHConfigAttrs.Comment,
						"trail_indent": SSHConfigAttrs.TrailIndent,
					},
				},
			},
		},
	},
}

var SSHConfigDataSourceSchema = dschema.Schema{
	MarkdownDescription: "Read and parse SSH config files",
	Attributes: map[string]dschema.Attribute{
		"path":    SSHConfigAttrs.Path,
		"content": SSHConfigAttrs.Content,
		"lines":   SSHConfigLineSchema,
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
		"lines":             SSHConfigLineSchema,
		"id":                SSHConfigAttrs.ID,
	},
}
