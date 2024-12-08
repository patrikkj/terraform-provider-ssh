provider "ssh" {
  # Connection using password
  host     = "internal-server.example.com"
  user     = "app-user"
  password = "supersecret"
  port     = 2222
}

provider "ssh" {
  # Direct connection using private key
  host        = "example.com"
  user        = "admin"
  private_key = file("~/.ssh/id_rsa")
  port        = 22
}

provider "ssh" {
  # Connection through bastion host
  host        = "internal-server.example.com"
  port        = 2222
  user        = "app-user"
  private_key = file("~/.ssh/app_key")

  bastion_host        = "bastion.example.com"
  bastion_port        = 22
  bastion_user        = "bastion-user"
  bastion_private_key = file("~/.ssh/bastion_key")
}
