package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/sshconf"
)

// Define the base attributes that will be used in both the main object and children
var sshConfigLineAttrTypes = map[string]attr.Type{
	"key":          types.StringType,
	"value":        types.StringType,
	"indent":       types.StringType,
	"sep":          types.StringType,
	"comment":      types.StringType,
	"trail_indent": types.StringType,
}

// Create the object type with the recursive children structure
var sshConfigLineObjectType = types.ObjectType{
	AttrTypes: func() map[string]attr.Type {
		attrs := make(map[string]attr.Type)
		for k, v := range sshConfigLineAttrTypes {
			attrs[k] = v
		}
		attrs["children"] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: sshConfigLineAttrTypes,
			},
		}
		return attrs
	}(),
}

func convertSSHConfLine(line sshconf.Line) SSHConfigLine {
	// Convert children to []attr.Value
	childrenValues := make([]attr.Value, len(line.Children))
	for i, child := range line.Children {
		converted := convertSSHConfLine(child)
		childrenValues[i] = converted.toAttrValue()
	}

	// Create types.List for children
	childrenList, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: sshConfigLineAttrTypes,
		},
		childrenValues,
	)

	return SSHConfigLine{
		Key:         types.StringValue(line.Key),
		Value:       types.StringValue(line.Value),
		Indent:      types.StringValue(line.Indent),
		Sep:         types.StringValue(line.Sep),
		Comment:     types.StringValue(line.Comment),
		TrailIndent: types.StringValue(line.TrailIndent),
		Children:    childrenList,
	}
}

// Helper method to convert SSHConfigLine to attr.Value
func (l SSHConfigLine) toAttrValue() attr.Value {
	return types.ObjectValueMust(
		sshConfigLineObjectType.AttrTypes,
		map[string]attr.Value{
			"key":          l.Key,
			"value":        l.Value,
			"indent":       l.Indent,
			"sep":          l.Sep,
			"comment":      l.Comment,
			"trail_indent": l.TrailIndent,
			"children":     l.Children,
		},
	)
}
