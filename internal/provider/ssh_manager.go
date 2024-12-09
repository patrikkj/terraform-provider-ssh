package provider

import (
	"fmt"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// SSHManager handles SSH connections for the provider
type SSHManager struct {
	providerClient *ssh.Client
}

// NewSSHManager creates a new SSH connection manager
func NewSSHManager(client *ssh.Client) *SSHManager {
	return &SSHManager{
		providerClient: client,
	}
}

// GetClient returns an SSH client based on the provided configuration
func (m *SSHManager) GetClient(config SSHConnectionConfig, useProviderAsBastion bool, bastion *SSHConnectionConfig, fromClient *ssh.Client) (*ssh.Client, bool, error) {
	// If useProviderAsBastion is true, recurse with the provider client
	if useProviderAsBastion {
		return m.GetClient(config, false, bastion, m.providerClient)
	}

	// If bastion is provided, get bastion client and recurse
	if bastion != nil {
		bastionClient, _, err := m.GetClient(*bastion, false, nil, fromClient)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to bastion: %v", err)
		}
		return m.GetClient(config, false, nil, bastionClient)
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
			return nil, false, fmt.Errorf("unable to parse private key: %v", err)
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
		// If config is not provided, use the provider client
		if *config.Host == "" {
			return m.providerClient, true, nil
		}
		client, err := ssh.Dial("tcp", target, sshConfig)
		if err != nil {
			return nil, false, fmt.Errorf("failed to connect to target host: %v", err)
		}
		return client, true, nil
	}

	// Create new client through fromClient
	conn, err := fromClient.Dial("tcp", target)
	if err != nil {
		return nil, false, fmt.Errorf("failed to connect to target host through bastion: %v", err)
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, target, sshConfig)
	if err != nil {
		return nil, false, fmt.Errorf("unable to create SSH connection through bastion: %v", err)
	}
	return ssh.NewClient(ncc, chans, reqs), true, nil
}
