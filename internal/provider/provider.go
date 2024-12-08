// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

var _ provider.Provider = &SSHProvider{}
var _ provider.ProviderWithFunctions = &SSHProvider{}

type SSHProvider struct {
	version string
}

type SSHProviderModel struct {
	Host              types.String `tfsdk:"host"`
	Port              types.Int64  `tfsdk:"port"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	PrivateKey        types.String `tfsdk:"private_key"`
	BastionHost       types.String `tfsdk:"bastion_host"`
	BastionPort       types.Int64  `tfsdk:"bastion_port"`
	BastionUser       types.String `tfsdk:"bastion_user"`
	BastionPassword   types.String `tfsdk:"bastion_password"`
	BastionPrivateKey types.String `tfsdk:"bastion_private_key"`
}

func (p *SSHProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ssh"
	resp.Version = p.version
}

func (p *SSHProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The hostname or IP address of the target SSH server",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the target SSH server",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username for SSH authentication",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for SSH authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "The private key for SSH authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"bastion_host": schema.StringAttribute{
				MarkdownDescription: "The hostname or IP address of the bastion host",
				Optional:            true,
			},
			"bastion_port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the bastion host",
				Optional:            true,
			},
			"bastion_user": schema.StringAttribute{
				MarkdownDescription: "The username for bastion host authentication",
				Optional:            true,
			},
			"bastion_password": schema.StringAttribute{
				MarkdownDescription: "The password for bastion host authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"bastion_private_key": schema.StringAttribute{
				MarkdownDescription: "The private key for bastion host authentication",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SSHProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SSHProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default port if not specified
	if config.Port.IsNull() {
		config.Port = types.Int64Value(22)
	}
	if config.BastionPort.IsNull() {
		config.BastionPort = types.Int64Value(22)
	}

	// Create SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            config.User.ValueString(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Add proper host key verification
	}

	// Configure authentication
	if !config.Password.IsNull() {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(config.Password.ValueString()))
	}
	if !config.PrivateKey.IsNull() {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey.ValueString()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse private key",
				err.Error(),
			)
			return
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	// Create SSH client
	client, err := createSSHClient(config, sshConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create SSH client",
			err.Error(),
		)
		return
	}

	// Store the client in the response
	resp.DataSourceData = client
	resp.ResourceData = client
}

func createSSHClient(config SSHProviderModel, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	target := net.JoinHostPort(config.Host.ValueString(), strconv.FormatInt(config.Port.ValueInt64(), 10))

	if config.BastionHost.IsNull() {
		// Direct connection
		return ssh.Dial("tcp", target, sshConfig)
	}

	// Connection through bastion host
	bastionConfig := &ssh.ClientConfig{
		User:            config.BastionUser.ValueString(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Add proper host key verification
	}

	if !config.BastionPassword.IsNull() {
		bastionConfig.Auth = append(bastionConfig.Auth, ssh.Password(config.BastionPassword.ValueString()))
	}
	if !config.BastionPrivateKey.IsNull() {
		signer, err := ssh.ParsePrivateKey([]byte(config.BastionPrivateKey.ValueString()))
		if err != nil {
			return nil, fmt.Errorf("unable to parse bastion private key: %v", err)
		}
		bastionConfig.Auth = append(bastionConfig.Auth, ssh.PublicKeys(signer))
	}

	bastionHost := net.JoinHostPort(config.BastionHost.ValueString(), strconv.FormatInt(config.BastionPort.ValueInt64(), 10))
	bastionClient, err := ssh.Dial("tcp", bastionHost, bastionConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to bastion host: %v", err)
	}

	conn, err := bastionClient.Dial("tcp", target)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to target host through bastion: %v", err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, target, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create SSH connection through bastion: %v", err)
	}

	return ssh.NewClient(ncc, chans, reqs), nil
}

func (p *SSHProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHExecResource,
		NewSSHFileResource,
	}
}

func (p *SSHProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSSHExecDataSource,
		NewSSHFileDataSource,
	}
}

func (p *SSHProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SSHProvider{
			version: version,
		}
	}
}
