# TODO ensure whitespace is being escaped
locals {
  AWS_ACCESS_KEY                  = trimspace(var.AWS_ACCESS_KEY)
  AWS_SECRET_KEY                  = trimspace(var.AWS_SECRET_KEY)
  AWS_IAM_PROFILE                 = trimspace(var.AWS_IAM_PROFILE)
  AWS_INSTANCE_TYPE               = trimspace(var.AWS_INSTANCE_TYPE)
  AWS_REGION                      = trimspace(var.AWS_REGION)
  AWS_EC2_AMI_NAME_FILTER         = trimspace(var.AWS_EC2_AMI_NAME_FILTER)
  AWS_EC2_AMI_ROOT_DEVICE_TYPE    = trimspace(var.AWS_EC2_AMI_ROOT_DEVICE_TYPE)
  AWS_EC2_AMI_VIRTUALIZATION_TYPE = trimspace(var.AWS_EC2_AMI_VIRTUALIZATION_TYPE)
  AWS_EC2_AMI_OWNERS              = distinct(var.AWS_EC2_AMI_OWNERS)
  AWS_USER_IDS                    = distinct(var.AWS_USER_IDS)
  AWS_AMI_NAME                    = lower(trimspace(regex_replace(var.AWS_AMI_NAME, "\\s+", "-")))
  AWS_EC2_SSH_USERNAME            = trimspace(var.AWS_EC2_SSH_USERNAME)
  AWS_EC2_INSTANCE_USERNAME       = lower(trimspace(regex_replace(var.AWS_EC2_INSTANCE_USERNAME, "\\s+", "-")))
  AWS_EC2_INSTANCE_USERNAME_HOME  = trimspace(var.AWS_EC2_INSTANCE_USERNAME_HOME)
  AWS_EC2_INSTANCE_SSH_KEY_NAME   = lower(trimspace(regex_replace(var.AWS_EC2_INSTANCE_SSH_KEY_NAME, "\\s+", "-")))

  AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL = trimspace(var.AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL)
  AWS_EC2_PUBLIC_DIRECTORY_INTERNAL          = trimspace(var.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL)

  GIT_USERNAME = lower(trimspace(regex_replace(var.GIT_USERNAME, "\\s+", "-")))
  GIT_EMAIL    = trimspace(var.GIT_EMAIL)

  ANSIBLE_TAGS = join(",", distinct(var.ANSIBLE_TAGS))
}

packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  access_key    = local.AWS_ACCESS_KEY
  secret_key    = local.AWS_SECRET_KEY
  profile       = local.AWS_IAM_PROFILE
  ami_name      = local.AWS_AMI_NAME
  instance_type = local.AWS_INSTANCE_TYPE
  region        = local.AWS_REGION
  ssh_username  = local.AWS_EC2_SSH_USERNAME

  source_ami_filter {
    filters = {
      name                = local.AWS_EC2_AMI_NAME_FILTER
      root-device-type    = local.AWS_EC2_AMI_ROOT_DEVICE_TYPE
      virtualization-type = local.AWS_EC2_AMI_VIRTUALIZATION_TYPE
    }
    most_recent = true
    owners      = local.AWS_EC2_AMI_OWNERS
  }

  temporary_security_group_source_public_ip = true
  ami_users                                 = local.AWS_USER_IDS

  tags = {
    Environment        = "development"
    Builder            = "packer"
    Deployer           = "pulumi"
    BuildRegion        = "{{ .BuildRegion }}"
    Base_AMI_ID        = "{{ .SourceAMI }}"
    Base_AMI_Name      = "{{ .SourceAMIName }}"
    Base_AMI_Owner     = "{{ .SourceAMIOwner }}"
    Base_AMI_OwnerName = "{{ .SourceAMIOwnerName }}"
  }
}

build {
  name = local.AWS_AMI_NAME

  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    env = {
      AWS_EC2_PUBLIC_DIRECTORY_INTERNAL : local.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL,
      AWS_EC2_SSH_USERNAME : local.AWS_EC2_SSH_USERNAME,
      AWS_EC2_PUBLIC_DIRECTORY_INTERNAL : local.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL,
    }
    scripts = [
      "./scripts/create_public_directory.sh",
      "./scripts/update_and_upgrade.sh",
      "./scripts/install_ansible.sh"
    ]
  }

  provisioner "ansible-local" {
    playbook_dir            = "./ansible"
    staging_directory       = local.AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL
    clean_staging_directory = true
    playbook_file           = "./ansible/playbooks/main.yml"
    extra_arguments = [
      "--tags",
      "${local.ANSIBLE_TAGS}",
      "--extra-vars",
      "AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL=${local.AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL}",
      "--extra-vars",
      "AWS_EC2_PUBLIC_DIRECTORY_INTERNAL=${local.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL}",
      "--extra-vars",
      "AWS_EC2_INSTANCE_USERNAME=${local.AWS_EC2_INSTANCE_USERNAME}",
      "--extra-vars",
      "AWS_EC2_INSTANCE_USERNAME_HOME=${local.AWS_EC2_INSTANCE_USERNAME_HOME}",
      "--extra-vars",
      "AWS_EC2_INSTANCE_SSH_KEY_NAME=${local.AWS_EC2_INSTANCE_SSH_KEY_NAME}",
      "--extra-vars",
      "AWS_EC2_SSH_USERNAME=${local.AWS_EC2_SSH_USERNAME}",
      "--extra-vars",
      "GIT_USERNAME=${local.GIT_USERNAME}",
      "--extra-vars",
      "GIT_EMAIL=${local.GIT_EMAIL}",
      "--extra-vars",
      "AWS_ACCESS_KEY=${local.AWS_ACCESS_KEY}",
      "--extra-vars",
      "AWS_SECRET_KEY=${local.AWS_SECRET_KEY}",
      "--extra-vars",
      "AWS_REGION=${local.AWS_REGION}"
    ]
  }

  provisioner "file" {
    source      = "${local.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL}/${local.AWS_EC2_INSTANCE_SSH_KEY_NAME}"
    destination = "../"
    direction   = "download"
  }

  provisioner "shell" {
    env = {
      AWS_EC2_PUBLIC_DIRECTORY_INTERNAL : local.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL,
    }
    scripts = ["./scripts/remove_public_directory.sh"]
  }
}



