package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHConfigDataSource(t *testing.T) {
	// Create a temporary SSH config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	configContent := `# Test SSH config
Host example
    HostName example.com
    User admin
    Port 2222

Host *
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* no-op */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read and check SSH config
			{
				Config: fmt.Sprintf(`
data "ssh_config" "test" {
  path = %q
}

output "first_host" {
  value = data.ssh_config.test.lines[1].value
}

output "first_host_hostname" {
  value = data.ssh_config.test.lines[2].value
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify basic attributes
					resource.TestCheckResourceAttr("data.ssh_config.test", "content", configContent),
					resource.TestCheckResourceAttrSet("data.ssh_config.test", "id"),

					// Verify first Host block
					resource.TestCheckResourceAttr("data.ssh_config.test", "lines.1.key", "Host"),
					resource.TestCheckResourceAttr("data.ssh_config.test", "lines.1.value", "example"),

					// Verify children
					resource.TestCheckResourceAttr("data.ssh_config.test", "lines.1.children.0.key", "HostName"),
					resource.TestCheckResourceAttr("data.ssh_config.test", "lines.1.children.0.value", "example.com"),
				),
			},
			// Test error handling with non-existent file
			{
				Config: `
provider "ssh" {
  host = "dummy-host.example.com"
  user = "dummy-user"
}

data "ssh_config" "test" {
  path = "/path/to/nonexistent/config"
}
`,
				ExpectError: regexp.MustCompile(`Failed to read SSH config file`),
			},
		},
	})
}
