package provider

import (
	"fmt"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// SSHManager handles SSH connections for the provider
type SSHManager struct {
	providerConfig  *SSHConnectionConfig
	providerBastion *SSHConnectionConfig
	providerClient  *ssh.Client
}

// NewSSHManager creates a new SSH connection manager
func NewSSHManager(config *SSHConnectionConfig, bastion *SSHConnectionConfig) (*SSHManager, error) {
	manager := &SSHManager{
		providerConfig:  config,
		providerBastion: bastion,
	}

	// Create initial provider client
	client, _, err := manager.GetClient(*config, false, bastion, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	manager.providerClient = client
	return manager, nil
}

// GetClient returns an SSH client based on the provided configuration
func (m *SSHManager) GetClient(config SSHConnectionConfig, useProviderAsBastion bool, bastion *SSHConnectionConfig, fromClient *ssh.Client) (*ssh.Client, bool, error) {
	// If useProviderAsBastion is true, use the provider client as bastion
	if useProviderAsBastion {
		return m.GetClient(config, false, bastion, m.providerClient)
	}

	// If bastion is provided, get bastion client and recurse
	if bastion != nil {
		bastionClient, _, err := m.GetClient(*bastion, false, nil, fromClient)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to bastion: %w", err)
		}
		return m.GetClient(config, false, nil, bastionClient)
	}

	// If there is on configuration, fall back to provider client
	if fromClient == nil && config.Host == nil {
		return m.providerClient, false, nil
	}

	// Create ssh client configuration
	sshConfig := &ssh.ClientConfig{
		User:            *config.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Configure authentication
	if config.Password != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(*config.Password))
	}
	if config.PrivateKey != nil {
		signer, err := ssh.ParsePrivateKey([]byte(*config.PrivateKey))
		if err != nil {
			return nil, false, fmt.Errorf("unable to parse private key: %w", err)
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	// Create target from port and host
	port := *config.Port
	if port == 0 {
		port = 22
	}
	target := net.JoinHostPort(*config.Host, strconv.FormatInt(port, 10))

	// If there is no fromClient, return a new client using ssh.Dial
	if fromClient == nil {
		client, err := ssh.Dial("tcp", target, sshConfig)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to target host: %w", err)
		}
		return client, true, nil
	}

	// Create new client through fromClient
	conn, err := fromClient.Dial("tcp", target)
	if err != nil {
		return nil, false, fmt.Errorf("failed to connect to target host through bastion: %w", err)
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, target, sshConfig)
	if err != nil {
		return nil, false, fmt.Errorf("unable to create SSH connection through bastion: %w", err)
	}
	return ssh.NewClient(ncc, chans, reqs), true, nil
}
