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
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	initialConfig := heredoc(`
        # Test SSH config
        Host existing
            HostName existing.com
            User admin
            Port 22

        Host *
            StrictHostKeyChecking no
            UserKnownHostsFile /dev/null
    `)
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: heredoc(fmt.Sprintf(`
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
                `, configPath)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),
					resource.TestCheckResourceAttr("ssh_config.test", "find", "Host example"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "true"),
					resource.TestCheckResourceAttrSet("ssh_config.test", "content"),
					resource.TestCheckResourceAttr("ssh_config.test", "lines.1.key", "Host"),
					resource.TestCheckResourceAttr("ssh_config.test", "lines.1.value", "existing"),
				),
			},
			{
				Config: heredoc(fmt.Sprintf(`
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
                `, configPath)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),
					resource.TestCheckResourceAttr("ssh_config.test", "find", "Host example"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "true"),
					resource.TestCheckResourceAttrSet("ssh_config.test", "content"),
				),
			},
		},
	})
}

func TestAccSSHConfigResource_NoDelete(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	initialConfig := heredoc(`
        # Test SSH config
        Host *
            StrictHostKeyChecking no
            UserKnownHostsFile /dev/null
    `)
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: heredoc(fmt.Sprintf(`
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
                `, configPath)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ssh_config.test", "id"),
					resource.TestCheckResourceAttr("ssh_config.test", "delete_on_destroy", "false"),
				),
			},
			{
				Config: "# Empty config to trigger destroy",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						content, err := os.ReadFile(configPath)
						if err != nil {
							return err
						}
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

func TestAccSSHConfigResource_SimpleMatch(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	initialConfig := heredoc(`
        Host example
            HostName old.example.com
            User olduser
    `)
	err := os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}

	expectedConfig := heredoc(`
        Host example
            HostName new.example.com
            User newuser
    `)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: heredoc(fmt.Sprintf(`
                    resource "ssh_config" "test" {
                        path = %q
                        find = "Host example"
                        patch = <<-EOT
                            Host example
                                HostName new.example.com
                                User newuser
                            EOT
                        delete_on_destroy = true
                    }
                `, configPath)),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["ssh_config.test"]
						if !ok {
							return fmt.Errorf("resource not found: ssh_config.test")
						}

						content := rs.Primary.Attributes["content"]
						if content != expectedConfig {
							return fmt.Errorf("expected config:\n%s\n\ngot:\n%s", expectedConfig, content)
						}

						fileContent, err := os.ReadFile(configPath)
						if err != nil {
							return fmt.Errorf("failed to read config file: %v", err)
						}

						if string(fileContent) != expectedConfig {
							return fmt.Errorf("file content mismatch.\nexpected:\n%s\n\ngot:\n%s", expectedConfig, string(fileContent))
						}

						return nil
					},
				),
			},
		},
	})
}
