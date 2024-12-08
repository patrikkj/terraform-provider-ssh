package provider

import (
	"fmt"
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
					resource.TestCheckResourceAttr("ssh_exec.basic", "output", "hello world\n"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("ssh_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("ssh_exec.nonzero_allowed", "exit_code", "1"),

					// Whoami command
					resource.TestCheckResourceAttr("ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("ssh_exec.whoami", "output", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("ssh_exec.whoami", "exit_code", "0"),

					// Multiline command
					resource.TestCheckResourceAttr("ssh_exec.multiline", "command", "echo \"Line 1\"\necho \"Line 2\"\n"),
					resource.TestCheckResourceAttr("ssh_exec.multiline", "exit_code", "0"),
					resource.TestCheckResourceAttr("ssh_exec.multiline", "output", "Line 1\nLine 2\n"),
				),
			},
			// Test updates to commands
			{
				Config: testAccSSHExecResourceConfigUpdates(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ssh_exec.basic", "command", "echo 'updated'"),
					resource.TestCheckResourceAttr("ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("ssh_exec.basic", "output", "updated\n"),
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

resource "ssh_exec" "on_destroy" {
  command = "echo 'hello world'"
  on_destroy = "echo 'on_destroy' > /tmp/on_destroy"
}

resource "ssh_exec" "nonzero_allowed" {
  command        = "false"
  fail_if_nonzero = false
}

resource "ssh_exec" "whoami" {
  command = "whoami"
}

resource "ssh_exec" "multiline" {
  command = <<-EOF
	  echo "Line 1"
	  echo "Line 2"
	EOF
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

resource "ssh_exec" "multiline" {
  command = <<-EOF
	  echo "Line 1"
	  echo "Line 2"
	EOF
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}
