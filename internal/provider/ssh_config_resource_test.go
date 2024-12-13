package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSSHConfigResource(t *testing.T) {
	// Create a temporary SSH config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	// Initial config content
	initialConfig := `# Test SSH config
Host existing
    HostName existing.com
    User admin
    Port 22

Host *
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
`
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* no-op */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: fmt.Sprintf(`
resource "ssh_config" "test" {
  path = %q
  find = "Host example"
  patch = <<-EOT
Host example
    HostName example.com
    User admin
    Port 2222
EOT
  delete_on_destroy = true
}

output "content" {
  value = ssh_config.test.content
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the resource exists
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),

					// Verify the patch was applied
					resource.TestCheckResourceAttr("ssh_config.test", "find", "Host example"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "true"),

					// Verify the content contains both the original and new config
					resource.TestCheckResourceAttrSet("ssh_config.test", "content"),
					resource.TestCheckResourceAttr("ssh_config.test", "lines.1.key", "Host"),
					resource.TestCheckResourceAttr("ssh_config.test", "lines.1.value", "existing"),
				),
			},
			// Update testing
			{
				Config: fmt.Sprintf(`
resource "ssh_config" "test" {
  path = %q
  find = "Host example"
  patch = <<-EOT
Host example
    HostName example.com
    User admin
    Port 3333
    IdentityFile ~/.ssh/example_key
EOT
  delete_on_destroy = true
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the resource still exists
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),

					// Verify the patch was updated
					resource.TestCheckResourceAttr("ssh_config.test", "find", "Host example"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "true"),

					// Verify the content was updated
					resource.TestCheckResourceAttrSet("ssh_config.test", "content"),
				),
			},
			// Import testing
			{
				ResourceName:      "ssh_config.test",
				ImportState:       true,
				ImportStateVerify: true,
				// These fields are not read back during import
				ImportStateVerifyIgnore: []string{"patch", "delete_on_destroy"},
			},
		},
	})
}

func TestAccSSHConfigResource_NoDelete(t *testing.T) {
	// Create a temporary SSH config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	// Initial config content
	initialConfig := `# Test SSH config
Host *
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
`
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* no-op */ },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with delete_on_destroy = false
			{
				Config: fmt.Sprintf(`
resource "ssh_config" "test" {
  path = %q
  find = "Host example"
  patch = <<-EOT
Host example
    HostName example.com
    User admin
    Port 2222
EOT
  delete_on_destroy = false
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "false"),
				),
			},
			// Verify the config remains after destroy
			{
				Config: fmt.Sprintf(`
# Empty config to trigger destroy of the previous resource
`),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						// Read the config file
						content, err := os.ReadFile(configPath)
						if err != nil {
							return err
						}

						// Verify the Host example section is still there
						if !strings.Contains(string(content), "Host example") {
							return fmt.Errorf("expected config to contain 'Host example' after destroy")
						}

						return nil
					},
				),
			},
		},
	})
}
