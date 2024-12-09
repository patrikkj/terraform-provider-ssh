package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Helper function to parse file mode
func parseFileMode(mode string) fs.FileMode {
	var result uint32
	if _, err := fmt.Sscanf(mode, "%o", &result); err != nil {
		return 0644 // Default to 0644 if parsing fails
	}
	return fs.FileMode(result)
}

// generateFileID creates a unique identifier for a file based on its path
func generateFileID(path string) string {
	h := md5.New()
	h.Write([]byte(path))
	return hex.EncodeToString(h.Sum(nil))
}

// readFile reads a file's contents over SFTP
func readFile(client *ssh.Client, path string) (string, error) {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	f, err := sftpClient.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read file contents: %w", err)
	}

	// Strip trailing newline when reading
	contentStr := string(content)
	if len(contentStr) > 0 && strings.HasSuffix(contentStr, "\n") {
		contentStr = contentStr[:len(contentStr)-1]
	}

	return contentStr, nil
}

// writeFile writes content to a file over SFTP
func writeFile(ctx context.Context, client *ssh.Client, path, content string, permissions string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Create directory if needed
	dirPath := filepath.Dir(path)
	if err := sftpClient.MkdirAll(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or truncate the file
	f, err := sftpClient.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Ensure content ends with newline
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Write content
	if _, err := f.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	// Close the file before changing permissions
	f.Close()

	// Try to set permissions, but don't fail if it doesn't work
	mode := parseFileMode(permissions)
	tflog.Debug(ctx, fmt.Sprintf("Attempting to chmod %s to %s", path, mode))

	if err := sftpClient.Chmod(path, mode); err != nil {
		// Log the permission error but don't fail the operation
		tflog.Warn(ctx, fmt.Sprintf("Warning: Could not set permissions on %s to %s: %s",
			path, mode, err))
	}

	return nil
}

// deleteFile deletes a file over SFTP
func deleteFile(client *ssh.Client, path string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	if err := sftpClient.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
