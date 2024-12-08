package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHExecDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHExecDataSourceConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "stdout", "hi\n"),

					// Whoami command
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "stdout", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "exit_code", "0"),

					// Non-zero exit with fail_if_nonzero = false
					resource.TestCheckResourceAttr("data.ssh_exec.nonzero_allowed", "command", "false"),
					resource.TestCheckResourceAttr("data.ssh_exec.nonzero_allowed", "exit_code", "1"),
					resource.TestCheckResourceAttr("data.ssh_exec.nonzero_allowed", "stdout", ""),
				),
			},
		},
	})
}

func testAccSSHExecDataSourceConfig(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

data "ssh_exec" "basic" {
  command = "echo hi"
}

data "ssh_exec" "whoami" {
  command = "whoami"
}

data "ssh_exec" "nonzero_allowed" {
  command = "false"
  fail_if_nonzero = false
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Test for expected failure when command returns non-zero with fail_if_nonzero = true
func TestAccSSHExecDataSource_FailIfNonZero(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSSHExecDataSourceConfigFailIfNonZero(t),
				ExpectError: regexp.MustCompile("Command exited with non-zero status"),
			},
		},
	})
}

func testAccSSHExecDataSourceConfigFailIfNonZero(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
  host     = "%s"
  user     = "%s"
  password = "%s"
}

data "ssh_exec" "nonzero_fail" {
  command = "false"
  fail_if_nonzero = true
}
`, getEnvVarOrSkip(t, "SSH_HOST"), getEnvVarOrSkip(t, "SSH_USER"), getEnvVarOrSkip(t, "SSH_PASSWORD"))
}

// Add this new test function
func TestAccSSHExecDataSource_PrivateKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckPrivateKey(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHExecDataSourceConfigPrivateKey(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic command execution
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "command", "echo hi"),
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "exit_code", "0"),
					resource.TestCheckResourceAttr("data.ssh_exec.basic", "stdout", "hi\n"),

					// Whoami command
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "command", "whoami"),
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "stdout", getEnvVarOrSkip(t, "SSH_USER")+"\n"),
					resource.TestCheckResourceAttr("data.ssh_exec.whoami", "exit_code", "0"),
				),
			},
		},
	})
}

func testAccSSHExecDataSourceConfigPrivateKey(t *testing.T) string {
	return fmt.Sprintf(`
provider "ssh" {
	host        = "%s"
	user        = "%s"
	private_key = file("%s")
}

data "ssh_exec" "basic" {
	command = "echo hi"
}

data "ssh_exec" "whoami" {
	command = "whoami"
}
`, getEnvVarOrSkip(t, "SSH_HOST"),
		getEnvVarOrSkip(t, "SSH_USER"),
		getEnvVarOrSkip(t, "SSH_PRIVATE_KEY_PATH"))
}
