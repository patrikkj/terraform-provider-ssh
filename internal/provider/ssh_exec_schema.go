package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

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
		"port":                    SSHConnectionSchema.Port,
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
		"port":                    SSHConnectionSchema.Port,
		"use_provider_as_bastion": SSHConnectionSchema.UseProviderAsBastion,
		"bastion":                 SSHConnectionSchema.Bastion,
	},
}
