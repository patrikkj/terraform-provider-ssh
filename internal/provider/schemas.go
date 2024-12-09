package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

var SSHConnectionSchema = struct {
	Host                 rschema.StringAttribute
	User                 rschema.StringAttribute
	Password             rschema.StringAttribute
	PrivateKey           rschema.StringAttribute
	UseProviderAsBastion rschema.BoolAttribute
}{
	Host: rschema.StringAttribute{
		MarkdownDescription: "Override the provider's host configuration",
		Optional:            true,
	},
	User: rschema.StringAttribute{
		MarkdownDescription: "Override the provider's user configuration",
		Optional:            true,
	},
	Password: rschema.StringAttribute{
		MarkdownDescription: "Override the provider's password configuration",
		Optional:            true,
		Sensitive:           true,
	},
	PrivateKey: rschema.StringAttribute{
		MarkdownDescription: "Override the provider's private key configuration",
		Optional:            true,
		Sensitive:           true,
	},
	UseProviderAsBastion: rschema.BoolAttribute{
		MarkdownDescription: "Use the provider's connection as a bastion host",
		Optional:            true,
	},
}

var SSHProviderAttrs = struct {
	Host              pschema.StringAttribute
	Port              pschema.Int64Attribute
	User              pschema.StringAttribute
	Password          pschema.StringAttribute
	PrivateKey        pschema.StringAttribute
	BastionHost       pschema.StringAttribute
	BastionPort       pschema.Int64Attribute
	BastionUser       pschema.StringAttribute
	BastionPassword   pschema.StringAttribute
	BastionPrivateKey pschema.StringAttribute
}{
	Host: pschema.StringAttribute{
		MarkdownDescription: "The hostname or IP address of the target SSH server",
		Required:            true,
	},
	Port: pschema.Int64Attribute{
		MarkdownDescription: "The port number of the target SSH server",
		Optional:            true,
	},
	User: pschema.StringAttribute{
		MarkdownDescription: "The username for SSH authentication",
		Required:            true,
	},
	Password: pschema.StringAttribute{
		MarkdownDescription: "The password for SSH authentication",
		Optional:            true,
		Sensitive:           true,
	},
	PrivateKey: pschema.StringAttribute{
		MarkdownDescription: "The private key for SSH authentication",
		Optional:            true,
		Sensitive:           true,
	},
	BastionHost: pschema.StringAttribute{
		MarkdownDescription: "The hostname or IP address of the bastion host",
		Optional:            true,
	},
	BastionPort: pschema.Int64Attribute{
		MarkdownDescription: "The port number of the bastion host",
		Optional:            true,
	},
	BastionUser: pschema.StringAttribute{
		MarkdownDescription: "The username for bastion host authentication",
		Optional:            true,
	},
	BastionPassword: pschema.StringAttribute{
		MarkdownDescription: "The password for bastion host authentication",
		Optional:            true,
		Sensitive:           true,
	},
	BastionPrivateKey: pschema.StringAttribute{
		MarkdownDescription: "The private key for bastion host authentication",
		Optional:            true,
		Sensitive:           true,
	},
}
var SSHExecAttrs = struct {
	Command                   rschema.StringAttribute
	Output                    rschema.StringAttribute
	ExitCode                  rschema.Int64Attribute
	FailIfNonzero__Resource   rschema.BoolAttribute
	FailIfNonzero__DataSource rschema.BoolAttribute
	OnDestroy                 rschema.StringAttribute
	ID                        rschema.StringAttribute
}{
	Command: rschema.StringAttribute{
		MarkdownDescription: "Command to execute",
		Required:            true,
	},
	Output: rschema.StringAttribute{
		MarkdownDescription: "Output of the command",
		Computed:            true,
	},
	ExitCode: rschema.Int64Attribute{
		MarkdownDescription: "Exit code of the command",
		Computed:            true,
	},
	FailIfNonzero__Resource: rschema.BoolAttribute{
		MarkdownDescription: "Whether to fail if the command returns a non-zero exit code. Defaults to true if not specified.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	},
	FailIfNonzero__DataSource: rschema.BoolAttribute{
		MarkdownDescription: "Whether to fail if the command returns a non-zero exit code",
		Optional:            true,
	},
	OnDestroy: rschema.StringAttribute{
		MarkdownDescription: "Command to execute when the resource is destroyed",
		Optional:            true,
	},
	ID: rschema.StringAttribute{
		MarkdownDescription: "Unique identifier for this execution",
		Computed:            true,
	},
}

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

var SSHProviderSchema = pschema.Schema{
	Attributes: map[string]pschema.Attribute{
		"host":                SSHProviderAttrs.Host,
		"port":                SSHProviderAttrs.Port,
		"user":                SSHProviderAttrs.User,
		"password":            SSHProviderAttrs.Password,
		"private_key":         SSHProviderAttrs.PrivateKey,
		"bastion_host":        SSHProviderAttrs.BastionHost,
		"bastion_port":        SSHProviderAttrs.BastionPort,
		"bastion_user":        SSHProviderAttrs.BastionUser,
		"bastion_password":    SSHProviderAttrs.BastionPassword,
		"bastion_private_key": SSHProviderAttrs.BastionPrivateKey,
	},
}

var SSHExecResourceSchema = rschema.Schema{
	MarkdownDescription: "Execute commands over SSH with potential side effects",
	Attributes: map[string]rschema.Attribute{
		"command":         SSHExecAttrs.Command,
		"output":          SSHExecAttrs.Output,
		"exit_code":       SSHExecAttrs.ExitCode,
		"fail_if_nonzero": SSHExecAttrs.FailIfNonzero__Resource,
		"on_destroy":      SSHExecAttrs.OnDestroy,
		"id":              SSHExecAttrs.ID,

		// Common SSH connection attributes
		"host":                    SSHConnectionSchema.Host,
		"user":                    SSHConnectionSchema.User,
		"password":                SSHConnectionSchema.Password,
		"private_key":             SSHConnectionSchema.PrivateKey,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
	},
}

var SSHExecDataSourceSchema = dschema.Schema{
	MarkdownDescription: "Execute commands over SSH",
	Attributes: map[string]dschema.Attribute{
		"command":         SSHExecAttrs.Command,
		"output":          SSHExecAttrs.Output,
		"exit_code":       SSHExecAttrs.ExitCode,
		"fail_if_nonzero": SSHExecAttrs.FailIfNonzero__DataSource,
		"id":              SSHExecAttrs.ID,

		// Common SSH connection attributes
		"host":                    SSHConnectionSchema.Host,
		"user":                    SSHConnectionSchema.User,
		"password":                SSHConnectionSchema.Password,
		"private_key":             SSHConnectionSchema.PrivateKey,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
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
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
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
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
	},
}
