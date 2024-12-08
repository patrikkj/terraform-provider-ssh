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
  fail_if_nonzero = true            # Optional: Fail if exit code is non-zero (default: true)
}
```

#### `ssh_file` - Read Files

```hcl
data "ssh_file" "example" {
  path          = "/path/to/file"   # Required: Remote file path
  fail_if_absent = true             # Optional: Fail if file doesn't exist
}
```

## Resources

#### `ssh_exec` - Execute Commands

```hcl
resource "ssh_exec" "example" {
  command         = "echo 'hello world'"    # Required: Command to execute
  fail_if_nonzero = true                    # Optional: Fail if exit code is non-zero (default: true)
}
```

#### `ssh_file` - Write Files

```hcl
resource "ssh_file" "example" {
  path        = "/path/to/file"     # Required: Remote file path
  content     = "file content"      # Required: File content
  permissions = "0644"              # Optional: File permissions (default: "0644")
}
```
