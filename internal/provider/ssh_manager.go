package provider

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

// SSHManager handles SSH connections for the provider
type SSHManager struct {
	providerClient *ssh.Client
}

// NewSSHManager creates a new SSH connection manager
func NewSSHManager(providerClient *ssh.Client) *SSHManager {
	return &SSHManager{
		providerClient: providerClient,
	}
}

type SSHConnectionConfig struct {
	Host                 types.String
	User                 types.String
	Password             types.String
	PrivateKey           types.String
	UseProviderAsBastion types.Bool
}

// GetClient returns an SSH client based on the provided configuration
func (m *SSHManager) GetClient(config *SSHConnectionConfig) (*ssh.Client, bool, error) {
	// If no override credentials are provided, return the provider's client
	if config.Host.IsNull() && config.User.IsNull() && config.Password.IsNull() && config.PrivateKey.IsNull() &&
		(config.UseProviderAsBastion.IsNull() || !config.UseProviderAsBastion.ValueBool()) {
		return m.providerClient, false, nil
	}

	// Create SSH config for the target
	sshConfig := &ssh.ClientConfig{
		User:            config.User.ValueString(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Configure authentication
	if !config.Password.IsNull() {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(config.Password.ValueString()))
	}
	if !config.PrivateKey.IsNull() {
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey.ValueString()))
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse private key: %w", err)
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	// If using provider as bastion
	if !config.UseProviderAsBastion.IsNull() && config.UseProviderAsBastion.ValueBool() {
		targetHost := net.JoinHostPort(config.Host.ValueString(), "22")

		conn, err := m.providerClient.Dial("tcp", targetHost)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect through bastion: %w", err)
		}

		ncc, chans, reqs, err := ssh.NewClientConn(conn, targetHost, sshConfig)
		if err != nil {
			return nil, false, fmt.Errorf("failed to create SSH connection through bastion: %w", err)
		}

		return ssh.NewClient(ncc, chans, reqs), true, nil
	}

	// Direct connection
	client, err := ssh.Dial("tcp", net.JoinHostPort(config.Host.ValueString(), "22"), sshConfig)
	return client, true, err
}

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
