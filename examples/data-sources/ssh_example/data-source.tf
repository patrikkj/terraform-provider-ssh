# Execute a command on the remote server
data "ssh_exec" "example" {
  command = "uname -a"
}

# Execute a command
data "ssh_exec" "example" {
  command = "echo hi"
}

# Execute a multiline script/command
data "ssh_exec" "multiline" {
  command = <<EOF
    #!/bin/bash
    echo "Starting script..."
    date
    uptime
    echo "Current directory:"
    pwd
    echo "Done!"
  EOF
}

# Execute a command with a specific working directory
data "ssh_exec" "with_dir" {
  command     = "pwd"
  working_dir = "/tmp"
}

# Read contents of a remote file
data "ssh_file" "example" {
  path = "/etc/hostname"
}

# Read a file with specific permissions check
data "ssh_file" "secure_file" {
  path           = "/etc/secret.conf"
  fail_if_absent = true
}
