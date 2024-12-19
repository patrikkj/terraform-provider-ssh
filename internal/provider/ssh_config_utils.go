package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrikkj/sshconf"
)

// Define the base attributes that will be used in both the main object and children
var sshConfigLineChildAttrTypes = map[string]attr.Type{
	"key":          types.StringType,
	"value":        types.StringType,
	"indent":       types.StringType,
	"sep":          types.StringType,
	"comment":      types.StringType,
	"trail_indent": types.StringType,
}

// Create the object type with the non-recursive children structure
var sshConfigLineObjectType = types.ObjectType{
	AttrTypes: func() map[string]attr.Type {
		attrs := make(map[string]attr.Type)
		for k, v := range sshConfigLineChildAttrTypes {
			attrs[k] = v
		}
		attrs["children"] = types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: sshConfigLineChildAttrTypes,
			},
		}
		return attrs
	}(),
}

func convertSSHConfLineNoChildren(line sshconf.LineNoChildren) SSHConfigLineChild {
	return SSHConfigLineChild{
		Key:         types.StringValue(line.Key),
		Value:       types.StringValue(line.Value),
		Indent:      types.StringValue(line.Indent),
		Sep:         types.StringValue(line.Sep),
		Comment:     types.StringValue(line.Comment),
		TrailIndent: types.StringValue(line.TrailIndent),
	}
}

func convertSSHConfLine(line sshconf.Line) SSHConfigLine {
	// Convert children to []attr.Value
	childrenValues := make([]attr.Value, len(line.Children))
	for i, child := range line.Children {
		converted := convertSSHConfLineNoChildren(child)
		childrenValues[i] = types.ObjectValueMust(
			sshConfigLineChildAttrTypes,
			map[string]attr.Value{
				"key":          converted.Key,
				"value":        converted.Value,
				"indent":       converted.Indent,
				"sep":          converted.Sep,
				"comment":      converted.Comment,
				"trail_indent": converted.TrailIndent,
			},
		)
	}

	// Create types.List for children
	childrenList, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: sshConfigLineChildAttrTypes,
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

func readOrCreateConfig(path string) (*sshconf.SSHConfig, error) {
	expandedPath := expandPath(path)
	config, err := sshconf.ParseConfigFile(expandedPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read SSH config file: %w", err)
		}
		// Create new empty config if file doesn't exist
		return sshconf.ParseConfig(""), nil
	}
	return config, nil
}

func updateModelFromConfig(config *sshconf.SSHConfig) (types.String, types.List, error) {
	content := config.Render()

	// Convert the lines to our model format
	lines := make([]SSHConfigLine, len(config.Lines()))
	for i, line := range config.Lines() {
		lines[i] = convertSSHConfLine(line)
	}

	lineValues := make([]attr.Value, len(lines))
	for i, line := range lines {
		lineValues[i] = line.toAttrValue()
	}

	linesList, diags := types.ListValue(
		sshConfigLineObjectType,
		lineValues,
	)
	if diags.HasError() {
		return types.String{}, types.List{}, fmt.Errorf("failed to create lines list: %v", diags)
	}

	return types.StringValue(content), linesList, nil
}
