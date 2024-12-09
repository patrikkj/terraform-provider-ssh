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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

var _ provider.Provider = &SSHProvider{}
var _ provider.ProviderWithFunctions = &SSHProvider{}

type SSHProvider struct {
	version string
	manager *SSHManager
}

func (p *SSHProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ssh"
	resp.Version = p.version
}

func (p *SSHProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = SSHProviderSchema
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
		keyContent := config.PrivateKey.ValueString()

		// Parse the private key
		signer, err := ssh.ParsePrivateKey([]byte(keyContent))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse private key",
				fmt.Sprintf("Private key parsing failed: %v", err),
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

	// Create the SSH manager
	p.manager = NewSSHManager(client)

	// Store both the client and manager in the response
	resp.DataSourceData = p.manager
	resp.ResourceData = p.manager
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
