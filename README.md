# Terraform SSH Provider

A Terraform provider for executing commands and managing files over SSH.

## Provider Configuration

```hcl
provider "ssh" {
  host = "app.example.com"             # Required: Target host address
  user = "admin"                       # Required: SSH username
  private_key = file("~/.ssh/id_rsa")  # Required: Private key authentication
  password = "your_password"           # Optional: Password authentication (alternative to private_key)
  port = 22                            # Optional: SSH port (default: 22)

  # Optional: Bastion configuration
  bastion = {
    host = "bastion.example.com"              # Required: Bastion host address
    user = "bastion-user"                     # Required: Bastion username
    private_key = file("~/.ssh/bastion_key")  # Required: Bastion authentication
    password = "bastion_password"             # Optional: Password authentication (alternative to private_key)
    port = 22                                 # Optional: Bastion port (default: 22)
  }
}

```

## Data Sources

#### `ssh_exec` - Execute Commands (read-only)

```hcl
data "ssh_exec" "example" {
  command = "systemctl status myapp"   # Required: Command to execute
  fail_if_nonzero = true               # Optional: Fail on non-zero exit

  # Optional: Connection overrides (same as resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true              # Use provider's bastion
  # bastion = { ... }                           # Custom bastion configuration

  depends_on = [ssh_exec.service_deployment]  # Optional: Data source dependencies
}

# Available outputs:
output "example" {
  value = {
    stdout     = data.ssh_exec.example.output    # The command's output
    exit_code  = data.ssh_exec.example.exit_code # The command's exit code
  }
}
```

#### `ssh_file` - Read Files

```hcl
data "ssh_file" "example" {
  path = "/etc/existing/config.yml"  # Required: Remote file path
  fail_if_absent = false             # Optional: Don't fail if file is missing

  # Optional: Connection overrides (same as resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true              # Use provider's bastion
  # bastion = { ... }                           # Custom bastion configuration
}

# Available outputs:
output "example" {
  value = {
    content = data.ssh_file.example.content      # The file's contents
    id      = data.ssh_file.example.id          # Unique identifier for this file
  }
}
```

## Resources

#### `ssh_exec` - Execute Commands

```hcl
resource "ssh_exec" "example" {
  # Required: Command to execute
  command = <<-EOT
    systemctl daemon-reload
    systemctl restart myapp
  EOT

  on_destroy = "systemctl stop myapp"  # Optional: Command to run on destruction
  fail_if_nonzero = true               # Optional: Fail on non-zero exit (defaults to true)

  # Optional: Connection overrides (same as file resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true              # Use provider's bastion
  # bastion = { ... }                           # Custom bastion configuration

  depends_on = [ssh_file.app_config]   # Optional: Resource dependencies
}

# Available outputs:
output "example" {
  value = {
    output     = ssh_exec.example.output    # The command's output
    exit_code  = ssh_exec.example.exit_code # The command's exit code
    id         = ssh_exec.example.id        # Unique identifier (same as command)
  }
}
```

#### `ssh_file` - Write Files

```hcl
resource "ssh_file" "example" {
  path = "/etc/myapp/config.json"  # Required: Remote file path

  # Required: File content
  content = jsonencode({
    database_url = "postgresql://db.internal:5432/myapp"
    api_key      = var.api_key
    environment  = var.environment
  })

  permissions = "0644"             # Optional: File permissions (defaults to "0644")
  delete_on_destroy = true         # Optional: Whether to delete on destroy (defaults to true)

  # Optional: Override provider connection settings
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # port = 2222                                 # Override provider port

  # use_provider_as_bastion = true              # Optional: Use provider's bastion as jump host

  # Optional: Custom bastion configuration (overrides provider's bastion)
  # bastion = {
  #   host        = "custom-jump.example.com"
  #   user        = "jump-user"
  #   private_key = file("~/.ssh/jump_key")
  #   port        = 22
  # }
}

# Available outputs:
output "example" {
  value = {
    content = data.ssh_file.example.content      # The file's contents
    id      = data.ssh_file.example.id           # Unique identifier for this file
  }
}
```

## Authentication

The provider supports two authentication methods:

1. Password authentication using the `password` attribute
2. Private key authentication using the `private_key` attribute

At least one authentication method must be provided. If both are provided, private key authentication will be attempted first.

## Bastion/Jump Host

For environments requiring a bastion (jump) host, configure the bastion-related attributes. The same authentication methods (password or private key) are supported for the bastion host.
