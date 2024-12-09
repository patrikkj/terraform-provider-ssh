package provider

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func generateExecID(command string, timestamp time.Time) string {
	h := md5.New()
	h.Write([]byte(command))
	h.Write([]byte(timestamp.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h.Sum(nil))
}

func executeSSHCommand(manager *SSHManager, config *SSHConnectionModel, command string, failIfNonzero bool) (string, int64, error) {
	client, newClient, err := manager.GetClient(config)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get SSH client: %w", err)
	}

	// If we created a new client, close it when done
	if newClient {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		return "", 0, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Execute command and capture output
	output, err := session.CombinedOutput(command)
	outputStr := string(output)

	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode := int64(exitErr.ExitStatus())
			if failIfNonzero && exitCode != 0 {
				return outputStr, exitCode, fmt.Errorf("command exited with non-zero status %d: %s", exitCode, outputStr)
			}
			return outputStr, exitCode, nil
		}
		return outputStr, -1, fmt.Errorf("failed to execute command: %w", err)
	}

	return outputStr, 0, nil
}
