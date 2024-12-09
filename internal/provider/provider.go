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

	// Create target address
	target := net.JoinHostPort(config.Host.ValueString(), strconv.FormatInt(config.Port.ValueInt64(), 10))

	// Create SSH client
	client, err := ssh.Dial("tcp", target, sshConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create SSH client",
			err.Error(),
		)
		return
	}

	// Create the SSH manager
	p.manager = NewSSHManager(client)
	resp.DataSourceData = p.manager
	resp.ResourceData = p.manager
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
