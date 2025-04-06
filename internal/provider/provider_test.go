package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ssh": providerserver.NewProtocol6WithError(New("test")()),
}

func init() {
	// Load the .env file from multiple possible locations
	envFiles := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			fmt.Printf("Successfully loaded env file: %s\n", file)
			break
		}
	}
}

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set for acceptance tests
	requiredEnvVars := []string{
		"SSH_HOST",
		"SSH_USER",
	}

	for _, envVar := range requiredEnvVars {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("Environment variable %s must be set for acceptance tests", envVar)
		}
	}

	// Optional environment variables - at least one authentication method must be provided
	if os.Getenv("SSH_PASSWORD") == "" && os.Getenv("SSH_PRIVATE_KEY_PATH") == "" {
		t.Fatal("Either SSH_PASSWORD or SSH_PRIVATE_KEY_PATH environment variable must be set for acceptance tests")
	}

	// Check bastion host configuration if provided
	if bastionHost := os.Getenv("SSH_BASTION_HOST"); bastionHost != "" {
		// If bastion host is specified, user must be provided
		if os.Getenv("SSH_BASTION_USER") == "" {
			t.Fatal("SSH_BASTION_USER must be set when SSH_BASTION_HOST is provided")
		}
		// If bastion host is specified, either password or private key must be provided
		if os.Getenv("SSH_BASTION_PASSWORD") == "" && os.Getenv("SSH_BASTION_PRIVATE_KEY") == "" {
			t.Fatal("Either SSH_BASTION_PASSWORD or SSH_BASTION_PRIVATE_KEY must be set when SSH_BASTION_HOST is provided")
		}
	}
}

func testAccPreCheckPrivateKey(t *testing.T) {
	requiredEnvVars := []string{
		"SSH_HOST",
		"SSH_USER",
		"SSH_PRIVATE_KEY_PATH",
	}

	for _, envVar := range requiredEnvVars {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("Environment variable %s must be set for private key authentication tests", envVar)
		}
	}
}
