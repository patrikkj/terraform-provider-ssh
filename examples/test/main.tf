terraform {
  required_version = ">= 1.0"

  required_providers {
    ssh = {
      source  = "local.providers/patrikkj/ssh"
      version = "~> 0.1.0"
    }
  }
}

provider "ssh" {}

resource "ssh_config" "test" {
  path  = "~/.ssh/config"
  find  = "Host example"
  patch = <<-EOT
    Host example
        HostName example.com
        User admin
        Port 2222
    EOT
}
