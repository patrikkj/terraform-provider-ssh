package provider

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

// SSHManager handles SSH connections for the provider
type SSHManager struct {
	Client *ssh.Client
}

// NewSSHManager creates a new SSH connection manager
func NewSSHManager(providerClient *ssh.Client) *SSHManager {
	return &SSHManager{
		Client: providerClient,
	}
}

// GetClient returns an SSH client based on the provided configuration
func (m *SSHManager) GetClient(config *SSHConnectionModel) (*ssh.Client, bool, error) {
	// If no override credentials are provided, return the provider's client
	if config.Host.IsNull() && config.User.IsNull() && config.Password.IsNull() && config.PrivateKey.IsNull() &&
		(config.UseProviderAsBastion.IsNull() || !config.UseProviderAsBastion.ValueBool()) {
		return m.Client, false, nil
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

		conn, err := m.Client.Dial("tcp", targetHost)
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
