package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHExecResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHExecResourceConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("ssh_exec.basic", "command", "echo 'hello world'"),
					resource.TestCheckResourceAttr("ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttrSet("ssh_exec.basic", "output"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("ssh_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("ssh_exec.nonzero_allowed", "exit_code", "1"),

					// Whoami command
					resource.TestCheckResourceAttr("ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("ssh_exec.whoami", "output", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("ssh_exec.whoami", "exit_code", "0"),

					// File write and verify
					resource.TestCheckResourceAttr("ssh_exec.file_write", "command",
						"echo 'test content' > /tmp/test_write.txt && ls -l /tmp/test_write.txt && cat /tmp/test_write.txt"),
					resource.TestCheckResourceAttr("ssh_exec.file_write", "exit_code", "0"),
					resource.TestMatchResourceAttr(
						"ssh_exec.file_write",
						"output",
						regexp.MustCompile(`-rw-.*\s+/tmp/test_write.txt\ntest content\n`),
					),
				),
			},
			// Test updates to commands
			{
				Config: testAccSSHExecResourceConfigUpdates(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ssh_exec.basic", "command", "echo 'updated'"),
					resource.TestCheckResourceAttr("ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("ssh_exec.basic", "output", "updated\n"),

					resource.TestCheckResourceAttr("ssh_exec.file_write", "command",
						"echo 'updated content' > /tmp/test_write.txt && cat /tmp/test_write.txt"),
					resource.TestCheckResourceAttr("ssh_exec.file_write", "exit_code", "0"),
					resource.TestCheckResourceAttr("ssh_exec.file_write", "output", "updated content\n"),
				),
			},
		},
	})
}

func testAccSSHExecResourceConfig(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

resource "ssh_exec" "basic" {
  command = "echo 'hello world'"
}

resource "ssh_exec" "nonzero_allowed" {
  command        = "false"
  fail_if_nonzero = false
}

resource "ssh_exec" "whoami" {
  command = "whoami"
}

resource "ssh_exec" "file_write" {
  command = "echo 'test content' > /tmp/test_write.txt && ls -l /tmp/test_write.txt && cat /tmp/test_write.txt"
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

func testAccSSHExecResourceConfigUpdates(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

resource "ssh_exec" "basic" {
  command = "echo 'updated'"
}

resource "ssh_exec" "nonzero_allowed" {
  command        = "false"
  fail_if_nonzero = false
}

resource "ssh_exec" "whoami" {
  command = "whoami"
}

resource "ssh_exec" "file_write" {
  command = "echo 'updated content' > /tmp/test_write.txt && cat /tmp/test_write.txt"
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}
