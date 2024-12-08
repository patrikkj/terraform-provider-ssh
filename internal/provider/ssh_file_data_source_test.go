package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHFileDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHFileDataSourceConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Test reading an existing file
					resource.TestCheckResourceAttr("data.ssh_file.hostname", "path", "/etc/hostname"),
					resource.TestCheckResourceAttrSet("data.ssh_file.hostname", "content"),

					// Test reading a non-existent file with fail_if_absent = false
					resource.TestCheckResourceAttr("data.ssh_file.missing_optional", "path", "/nonexistent/file"),
					resource.TestCheckResourceAttr("data.ssh_file.missing_optional", "content", ""),

					// Test reading /etc/hosts file
					resource.TestCheckResourceAttr("data.ssh_file.hosts", "path", "/etc/hosts"),
					resource.TestCheckResourceAttrSet("data.ssh_file.hosts", "content"),
				),
			},
		},
	})
}

func testAccSSHFileDataSourceConfig(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

data "ssh_file" "hostname" {
  path = "/etc/hostname"
}

data "ssh_file" "missing_optional" {
  path = "/nonexistent/file"
  fail_if_absent = false
}

data "ssh_file" "hosts" {
  path = "/etc/hosts"
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Test for expected failure when reading non-existent file with fail_if_absent = true
func TestAccSSHFileDataSource_FailIfAbsent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSSHFileDataSourceConfigFailIfAbsent(t),
				ExpectError: regexp.MustCompile("Failed to read file"),
			},
		},
	})
}

func testAccSSHFileDataSourceConfigFailIfAbsent(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

data "ssh_file" "missing_required" {
  path = "/nonexistent/file"
  fail_if_absent = true
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Helper function to get environment variables or skip test
func getEnvVarOrSkip(t *testing.T, name string) string {
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("Environment variable %s is not set", name)
	}
	return value
}
