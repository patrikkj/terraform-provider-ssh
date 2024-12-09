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
func (m *SSHManager) GetClient(config *SSHConnectionModel) (*ssh.Client, bool, error) {
	// If no override configuration is provided, return the provider client
	if config.Host.IsNull() && !config.UseProviderAsBastion.ValueBool() {
		return m.providerClient, false, nil
	}

	// Create SSH client configuration
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
			return nil, false, fmt.Errorf("unable to parse private key: %v", err)
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	// Handle bastion configuration
	var startingClient *ssh.Client
	if config.UseProviderAsBastion.ValueBool() {
		startingClient = m.providerClient
	}

	client, err := m.createSSHClientWithBastion(config, sshConfig, startingClient)
	if err != nil {
		return nil, false, err
	}

	return client, true, nil
}

func (m *SSHManager) createSSHClientWithBastion(config *SSHConnectionModel, sshConfig *ssh.ClientConfig, startingClient *ssh.Client) (*ssh.Client, error) {
	target := net.JoinHostPort(config.Host.ValueString(), "22") // Default to port 22 if not specified

	// If no bastion is configured and we're not using provider as bastion, make direct connection
	if config.BastionHost.IsNull() && startingClient == nil {
		return ssh.Dial("tcp", target, sshConfig)
	}

	var bastionClient *ssh.Client
	if startingClient != nil {
		bastionClient = startingClient
	} else {
		// Configure bastion connection
		bastionConfig := &ssh.ClientConfig{
			User:            config.BastionUser.ValueString(),
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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

		bastionPort := "22"
		if !config.BastionPort.IsNull() {
			bastionPort = strconv.FormatInt(config.BastionPort.ValueInt64(), 10)
		}
		bastionHost := net.JoinHostPort(config.BastionHost.ValueString(), bastionPort)

		var err error
		bastionClient, err = ssh.Dial("tcp", bastionHost, bastionConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to bastion host: %v", err)
		}
	}

	// Connect to target through bastion
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
