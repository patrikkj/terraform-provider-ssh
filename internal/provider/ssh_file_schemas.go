package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

var SSHFileAttrs = struct {
	Path                rschema.StringAttribute
	Content__Resource   rschema.StringAttribute
	Content__DataSource rschema.StringAttribute
	Permissions         rschema.StringAttribute
	FailIfAbsent        rschema.BoolAttribute
	DeleteOnDestroy     rschema.BoolAttribute
	ID                  rschema.StringAttribute
}{
	Path: rschema.StringAttribute{
		MarkdownDescription: "Path to the file",
		Required:            true,
	},
	Content__Resource: rschema.StringAttribute{
		MarkdownDescription: "Content of the file",
		Required:            true,
	},
	Content__DataSource: rschema.StringAttribute{
		MarkdownDescription: "Content of the file",
		Computed:            true,
	},
	Permissions: rschema.StringAttribute{
		MarkdownDescription: "File permissions (e.g., '0644')",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("0644"),
	},
	FailIfAbsent: rschema.BoolAttribute{
		MarkdownDescription: "Whether to fail if the file does not exist",
		Optional:            true,
	},
	DeleteOnDestroy: rschema.BoolAttribute{
		MarkdownDescription: "Whether to delete the file when the resource is destroyed. Defaults to true.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
	ID: rschema.StringAttribute{
		MarkdownDescription: "Unique identifier for this file",
		Computed:            true,
	},
}

var SSHFileResourceSchema = rschema.Schema{
	MarkdownDescription: "Manage files over SSH with potential side effects",
	Attributes: map[string]rschema.Attribute{
		"path":              SSHFileAttrs.Path,
		"content":           SSHFileAttrs.Content__Resource,
		"permissions":       SSHFileAttrs.Permissions,
		"fail_if_absent":    SSHFileAttrs.FailIfAbsent,
		"delete_on_destroy": SSHFileAttrs.DeleteOnDestroy,
		"id":                SSHFileAttrs.ID,

		// Common SSH connection attributes
		"host":                    SSHConnectionSchema.Host,
		"user":                    SSHConnectionSchema.User,
		"password":                SSHConnectionSchema.Password,
		"private_key":             SSHConnectionSchema.PrivateKey,
		"port":                    SSHConnectionSchema.Port,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
		"bastion":                 SSHConnectionSchema.Bastion,
	},
}

var SSHFileDataSourceSchema = dschema.Schema{
	MarkdownDescription: "Read files over SSH",
	Attributes: map[string]dschema.Attribute{
		"path":           SSHFileAttrs.Path,
		"content":        SSHFileAttrs.Content__DataSource,
		"permissions":    SSHFileAttrs.Permissions,
		"fail_if_absent": SSHFileAttrs.FailIfAbsent,
		"id":             SSHFileAttrs.ID,

		// Common SSH connection attributes
		"host":                    SSHConnectionSchema.Host,
		"user":                    SSHConnectionSchema.User,
		"password":                SSHConnectionSchema.Password,
		"private_key":             SSHConnectionSchema.PrivateKey,
		"port":                    SSHConnectionSchema.Port,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
		"bastion":                 SSHConnectionSchema.Bastion,
	},
}
