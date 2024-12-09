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
	Port                 rschema.Int64Attribute
	UseProviderAsBastion rschema.BoolAttribute
	Bastion              rschema.SingleNestedAttribute
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
	Port: rschema.Int64Attribute{
		MarkdownDescription: "The port number to connect to",
		Optional:            true,
	},
	UseProviderAsBastion: rschema.BoolAttribute{
		MarkdownDescription: "Use the provider's connection as a bastion host",
		Optional:            true,
	},
	Bastion: rschema.SingleNestedAttribute{
		MarkdownDescription: "Bastion host configuration",
		Optional:            true,
		Attributes: map[string]rschema.Attribute{
			"host": rschema.StringAttribute{
				MarkdownDescription: "The hostname or IP address of the bastion host",
				Required:            true,
			},
			"port": rschema.Int64Attribute{
				MarkdownDescription: "The port number of the bastion host",
				Optional:            true,
			},
			"user": rschema.StringAttribute{
				MarkdownDescription: "The username for bastion host authentication",
				Required:            true,
			},
			"password": rschema.StringAttribute{
				MarkdownDescription: "The password for bastion host authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"private_key": rschema.StringAttribute{
				MarkdownDescription: "The private key for bastion host authentication",
				Optional:            true,
				Sensitive:           true,
			},
		},
	},
}

var SSHBastionAttr = rschema.SingleNestedAttribute{
	MarkdownDescription: "Bastion host configuration",
	Optional:            true,
	Attributes: map[string]rschema.Attribute{
		"host":        SSHConnectionSchema.Host,
		"port":        SSHConnectionSchema.Port,
		"user":        SSHConnectionSchema.User,
		"password":    SSHConnectionSchema.Password,
		"private_key": SSHConnectionSchema.PrivateKey,
	},
}

var SSHProviderAttrs = struct {
	Host       rschema.StringAttribute
	Port       rschema.Int64Attribute
	User       rschema.StringAttribute
	Password   rschema.StringAttribute
	PrivateKey rschema.StringAttribute
	Bastion    rschema.SingleNestedAttribute
}{
	Host: rschema.StringAttribute{
		MarkdownDescription: "The hostname or IP address of the target SSH server",
		Required:            true,
	},
	Port: rschema.Int64Attribute{
		MarkdownDescription: "The port number of the target SSH server",
		Optional:            true,
	},
	User: rschema.StringAttribute{
		MarkdownDescription: "The username for SSH authentication",
		Required:            true,
	},
	Password: rschema.StringAttribute{
		MarkdownDescription: "The password for SSH authentication",
		Optional:            true,
		Sensitive:           true,
	},
	PrivateKey: rschema.StringAttribute{
		MarkdownDescription: "The private key for SSH authentication",
		Optional:            true,
		Sensitive:           true,
	},
	Bastion: SSHBastionAttr,
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
		"host":        SSHProviderAttrs.Host,
		"port":        SSHProviderAttrs.Port,
		"user":        SSHProviderAttrs.User,
		"password":    SSHProviderAttrs.Password,
		"private_key": SSHProviderAttrs.PrivateKey,
		"bastion":     SSHProviderAttrs.Bastion,
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
		"bastion":                 SSHConnectionSchema.Bastion,
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
		"bastion":                 SSHConnectionSchema.Bastion,
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
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
		"bastion":                 SSHConnectionSchema.Bastion,
	},
}
