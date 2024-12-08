package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// GetCommonSSHConnectionSchema returns schema attributes common to all resources/data sources
func GetCommonSSHConnectionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"host": schema.StringAttribute{
			MarkdownDescription: "Override the provider's host configuration",
			Optional:            true,
		},
		"user": schema.StringAttribute{
			MarkdownDescription: "Override the provider's user configuration",
			Optional:            true,
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "Override the provider's password configuration",
			Optional:            true,
			Sensitive:           true,
		},
		"private_key": schema.StringAttribute{
			MarkdownDescription: "Override the provider's private key configuration",
			Optional:            true,
			Sensitive:           true,
		},
		"use_provider_as_bastion": schema.BoolAttribute{
			MarkdownDescription: "Use the provider's connection as a bastion host",
			Optional:            true,
		},
	}
}

// GetSSHExecSchema returns the schema for SSH exec resources/data sources
func GetSSHExecSchema() map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"command": schema.StringAttribute{
			MarkdownDescription: "Command to execute",
			Required:            true,
		},
		"output": schema.StringAttribute{
			MarkdownDescription: "Output of the command",
			Computed:            true,
		},
		"exit_code": schema.Int64Attribute{
			MarkdownDescription: "Exit code of the command",
			Computed:            true,
		},
		"fail_if_nonzero": schema.BoolAttribute{
			MarkdownDescription: "Whether to fail if the command returns a non-zero exit code. Defaults to true if not specified.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"on_destroy": schema.StringAttribute{
			MarkdownDescription: "Command to execute when the resource is destroyed",
			Optional:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this execution",
			Computed:            true,
		},
	}

	// Merge with common SSH connection attributes
	for k, v := range GetCommonSSHConnectionSchema() {
		attributes[k] = v
	}

	return attributes
}

// GetSSHFileSchema returns the schema for SSH file resources/data sources
func GetSSHFileSchema() map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"path": schema.StringAttribute{
			MarkdownDescription: "Path to the file",
			Required:            true,
		},
		"content": schema.StringAttribute{
			MarkdownDescription: "Content of the file",
			Required:            true,
		},
		"permissions": schema.StringAttribute{
			MarkdownDescription: "File permissions (e.g., '0644')",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("0644"),
		},
		"fail_if_absent": schema.BoolAttribute{
			MarkdownDescription: "Whether to fail if the file does not exist",
			Optional:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this file",
			Computed:            true,
		},
	}

	// Merge with common SSH connection attributes
	for k, v := range GetCommonSSHConnectionSchema() {
		attributes[k] = v
	}

	return attributes
}
