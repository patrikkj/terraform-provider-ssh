# Provider Configuration
# ---------------------
provider "ssh" {
  host = "app.example.com" # Required: Target host address
  user = "admin"           # Required: SSH username

  # Authentication - Either password or private_key is required
  private_key = file("~/.ssh/id_rsa")
  # password = "your_password"    # Alternative to private_key

  port = 22 # Optional: SSH port (defaults to 22)

  # Optional: Bastion/jump host configuration
  bastion = {
    host = "bastion.example.com" # Required if bastion block is present
    user = "bastion-user"        # Required if bastion block is present

    # Authentication - Either password or private_key is required
    private_key = file("~/.ssh/bastion_key")
    # password = "bastion_password"

    port = 22 # Optional: Bastion port (defaults to 22)
  }
}

# File Resource Examples
# --------------------

# Basic file management
resource "ssh_file" "app_config" {
  path = "/etc/myapp/config.json" # Required: Remote file path

  # Required: File content
  content = jsonencode({
    database_url = "postgresql://db.internal:5432/myapp"
    api_key      = var.api_key
    environment  = var.environment
  })

  permissions       = "0644" # Optional: File permissions (defaults to "0644")
  delete_on_destroy = true   # Optional: Whether to delete on destroy (defaults to true)

  # Optional: Override provider connection settings
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # port = 2222                                 # Override provider port

  # use_provider_as_bastion = true             # Optional: Use provider's bastion as jump host

  # Optional: Custom bastion configuration (overrides provider's bastion)
  # bastion = {
  #   host        = "custom-jump.example.com"
  #   user        = "jump-user"
  #   private_key = file("~/.ssh/jump_key")
  #   port        = 22
  # }
}

# File Data Source Example
# ----------------------
data "ssh_file" "existing_config" {
  path           = "/etc/existing/config.yml" # Required: Remote file path
  fail_if_absent = false                      # Optional: Don't fail if file is missing

  # Optional: Connection overrides (same as resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true             # Use provider's bastion
  # bastion = { ... }                          # Custom bastion configuration
}

# Command Execution Resource Example
# -------------------------------
resource "ssh_exec" "service_deployment" {
  # Required: Command to execute
  command = <<-EOT
    systemctl daemon-reload
    systemctl restart myapp
  EOT

  on_destroy      = "systemctl stop myapp" # Optional: Command to run on destruction
  fail_if_nonzero = true                   # Optional: Fail on non-zero exit (defaults to true)

  # Optional: Connection overrides (same as file resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true             # Use provider's bastion
  # bastion = { ... }                          # Custom bastion configuration

  depends_on = [ssh_file.app_config] # Optional: Resource dependencies
}

# Command Execution Data Source Example
# ---------------------------------
data "ssh_exec" "service_status" {
  command         = "systemctl status myapp" # Required: Command to execute
  fail_if_nonzero = true                     # Optional: Fail on non-zero exit

  # Optional: Connection overrides (same as resource)
  # host = "different-host.example.com"         # Override provider host
  # user = "different-user"                     # Override provider user
  # private_key = file("~/.ssh/different_key")  # Override provider authentication
  # use_provider_as_bastion = true             # Use provider's bastion
  # bastion = { ... }                          # Custom bastion configuration

  depends_on = [ssh_exec.service_deployment] # Optional: Data source dependencies
}

# Outputs
# -------
output "config_content" {
  value     = data.ssh_file.existing_config.content
  sensitive = true # Mark as sensitive if content contains secrets
}

output "service_status" {
  value = data.ssh_exec.service_status.output
}

# Variables
# --------
variable "api_key" {
  type      = string
  sensitive = true # Mark as sensitive since it's a secret
}

variable "environment" {
  type    = string
  default = "production" # Default environment value
}
