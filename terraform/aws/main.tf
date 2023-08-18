locals {
  AWS_ACCESS_KEY_ID = trimspace(var.AWS_ACCESS_KEY_ID)
  AWS_SECRET_ACCESS_KEY = trimspace(var.AWS_SECRET_ACCESS_KEY)
  AWS_REGION                   = trimspace(var.AWS_REGION)
  AMI_ID                       = trimspace(var.AMI_ID)
  AWS_INSTANCE_TYPE            = trimspace(var.AWS_INSTANCE_TYPE)
  AWS_EC2_INSTANCE_VOLUME_TYPE = trimspace(var.AWS_EC2_INSTANCE_VOLUME_TYPE)
  AWS_EC2_INSTANCE_VOLUME_SIZE = var.AWS_EC2_INSTANCE_VOLUME_SIZE

  ALLOWED_IP   = trimspace(var.ALLOWED_IP)
  KEY_NAME     = trimspace(var.KEY_NAME)
  SSH_PORT     = var.SSH_PORT
  IP_FILE_NAME = trimspace(var.IP_FILE_NAME)
  USERNAME     = trimspace(var.USERNAME)

  first_available_zone = length(data.aws_availability_zones.available.names) > 0 ? data.aws_availability_zones.available.names[0] : null
  KUMO_NAME_TAG       = "kumo"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = local.AWS_REGION
  access_key = local.AWS_ACCESS_KEY_ID
  secret_key = local.AWS_SECRET_ACCESS_KEY
}

data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_ami" "kumo-ami" {
  most_recent = true

  filter {
    name   = "image-id"
    values = [local.AMI_ID]
  }
}

resource "aws_vpc" "kumo-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "aws_internet_gateway" "kumo-internet-gateway" {
  vpc_id = aws_vpc.kumo-vpc.id

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "aws_subnet" "kumo-subnet" {
  vpc_id                  = aws_vpc.kumo-vpc.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "aws_route_table" "kumo-route-table" {
  vpc_id = aws_vpc.kumo-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.kumo-internet-gateway.id
  }

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "aws_route_table_association" "kumo-route-table-association" {
  route_table_id = aws_route_table.kumo-route-table.id
  subnet_id      = aws_subnet.kumo-subnet.id
}

resource "aws_security_group" "kumo-security-group" {
  name   = "kumo-security-group"
  vpc_id = aws_vpc.kumo-vpc.id
  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "aws_vpc_security_group_egress_rule" "kumo-security-group-egress-rule" {
  security_group_id = aws_security_group.kumo-security-group.id

  cidr_ipv4   = "0.0.0.0/0"
  from_port   = 0
  ip_protocol = "-1"
  to_port     = 0
}

resource "aws_vpc_security_group_ingress_rule" "kumo-security-group-ingress-rule" {
  security_group_id = aws_security_group.kumo-security-group.id

  cidr_ipv4   = local.ALLOWED_IP
  from_port   = local.SSH_PORT
  ip_protocol = "tcp"
  to_port     = local.SSH_PORT
}

resource "tls_private_key" "kumo-ssh-key" {
  algorithm = "ED25519"
}

resource "aws_key_pair" "kumo-ssh-key-pair" {
  key_name   = local.KEY_NAME
  public_key = tls_private_key.kumo-ssh-key.public_key_openssh

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "local_file" "kumo-ssh-private-key" {
  content  = tls_private_key.kumo-ssh-key.private_key_openssh
  file_permission = "0600"
  filename = local.KEY_NAME
}

resource "aws_instance" "kumo-ec2-instance" {
  instance_type = local.AWS_INSTANCE_TYPE
  ami           = data.aws_ami.kumo-ami.id
  vpc_security_group_ids = [
    aws_security_group.kumo-security-group.id
  ]
  subnet_id         = aws_subnet.kumo-subnet.id
  availability_zone = local.first_available_zone
  key_name          = aws_key_pair.kumo-ssh-key-pair.key_name

  root_block_device {
    volume_type = local.AWS_EC2_INSTANCE_VOLUME_TYPE
    volume_size = local.AWS_EC2_INSTANCE_VOLUME_SIZE
  }

  user_data = <<-EOF
    #!/bin/bash
    path=/home/${local.USERNAME}/.ssh/authorized_keys
    user=${local.USERNAME}

    # Create the .ssh directory and authorized_keys file if they don't exist
    mkdir -p "/home/${local.USERNAME}/.ssh" && touch $path

    # Add SSH key to the authorized_keys file of a user
    echo "${tls_private_key.kumo-ssh-key.public_key_openssh}" >> $path

    # Set correct permissions for the authorized_keys file
    chmod 600 $path
    chown $user:$user $path
  EOF

  tags = {
    Name = local.KUMO_NAME_TAG
  }
}

resource "local_file" "kumo-ec2-public-ip" {
  content  = aws_instance.kumo-ec2-instance.public_ip
  filename = local.IP_FILE_NAME
}

output "connection_cmd" {
  description = "How to connect to your instance"
  value       = "Deployment complete, now try 'ssh -F kumossh kumo' to connect to your instance"
}