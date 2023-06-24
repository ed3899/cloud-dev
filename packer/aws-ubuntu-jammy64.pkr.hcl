packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  access_key    = var.AWS_ACCESS_KEY
  secret_key    = var.AWS_SECRET_KEY
  profile       = var.AWS_IAM_PROFILE
  ami_name      = var.AWS_AMI_NAME
  instance_type = var.AWS_INSTANCE_TYPE
  region        = var.AWS_REGION
  ssh_username  = var.AWS_INSTANCE_SSH_USERNAME

  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }

  ami_users = [var.AWS_USER_ID]
}

build {
  name = "cloud_dev"

  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    script = "./scripts/update_and_upgrade.sh"
  }

  provisioner "shell" {
    script = "./scripts/install_ansible.sh"
  }

  provisioner "ansible-local" {
    playbook_dir    = "./ansible"
    staging_directory = "/tmp/ansible"
    clean_staging_directory = true
    playbook_file = "./ansible/playbooks/main.yml"
    extra_arguments = ["--tags", "always"]
  }

  provisioner "file" {
    source = ""
    destination = ""
    direction = "download"
  }
}



