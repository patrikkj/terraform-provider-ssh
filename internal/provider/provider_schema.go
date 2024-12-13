package provider

import (
	pschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		Optional:            true,
	},
	Port: rschema.Int64Attribute{
		MarkdownDescription: "The port number of the target SSH server",
		Optional:            true,
	},
	User: rschema.StringAttribute{
		MarkdownDescription: "The username for SSH authentication",
		Optional:            true,
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
