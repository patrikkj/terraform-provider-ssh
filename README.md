# Terraform SSH Provider

A Terraform provider for executing commands and managing files over SSH.

## Provider Configuration

```hcl
provider "ssh" {
  host              = "example.com"      # Required: SSH host
  port              = 22                 # Optional: SSH port (default: 22)
  user              = "username"         # Required: SSH username
  password          = "password"         # Optional: SSH password
  private_key       = "key_content"      # Optional: SSH private key content

  # Optional: Bastion/jump host configuration
  bastion_host       = "bastion.example.com"
  bastion_port       = 22
  bastion_user       = "bastion_user"
  bastion_password   = "bastion_password"
  bastion_private_key = "bastion_key_content"
}
```

## Data Sources

#### `ssh_exec` - Execute Commands (read-only)

```hcl
data "ssh_exec" "example" {
  command         = "whoami"        # Required: Command to execute
  fail_if_nonzero = true           # Optional: Fail if exit code is non-zero (default: true)
}

# Available outputs:
output "command_result" {
  value = {
    stdout     = data.ssh_exec.example.output    # The command's output
    exit_code  = data.ssh_exec.example.exit_code # The command's exit code
  }
}
```

#### `ssh_file` - Read Files

```hcl
data "ssh_file" "example" {
  path          = "/path/to/file"   # Required: Remote file path
  fail_if_absent = true             # Optional: Fail if file doesn't exist
}

# Available outputs:
output "file_content" {
  value = data.ssh_file.example.content  # The file's content
}
```

## Resources

#### `ssh_exec` - Execute Commands

```hcl
resource "ssh_exec" "example" {
  command         = "echo 'hello world'"    # Required: Command to execute
  fail_if_nonzero = true                   # Optional: Fail if exit code is non-zero (default: true)
}

# Available outputs:
output "command_result" {
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
  path        = "/path/to/file"     # Required: Remote file path
  content     = "file content"      # Required: File content
  permissions = "0644"              # Optional: File permissions (default: "0644")
}

# Available outputs:
output "file_info" {
  value = {
    id          = ssh_file.example.id          # Unique identifier (same as path)
    path        = ssh_file.example.path        # Path to the file
    permissions = ssh_file.example.permissions # File permissions
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
