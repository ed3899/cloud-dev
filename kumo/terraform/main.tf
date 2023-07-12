locals {
  AWS_REGION                   = trimspace(var.AWS_REGION)
  AMI_ID                       = trimspace(var.AMI_ID)
  AWS_INSTANCE_TYPE            = trimspace(var.AWS_INSTANCE_TYPE)
  AWS_EC2_INSTANCE_VOLUME_TYPE = trimspace(var.AWS_EC2_INSTANCE_VOLUME_TYPE)
  AWS_EC2_INSTANCE_VOLUME_SIZE = var.AWS_EC2_INSTANCE_VOLUME_SIZE

  allowed_ip           = trimspace(var.allowed_ip)
  first_available_zone = length(data.aws_availability_zones.available.names) > 0 ? data.aws_availability_zones.available.names[0] : null
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
}

resource "aws_internet_gateway" "kumo-internet-gateway" {
  vpc_id     = aws_vpc.kumo-vpc.id
  depends_on = [aws_vpc.kumo-vpc]
}

resource "aws_subnet" "kumo-subnet" {
  vpc_id                  = aws_vpc.kumo-vpc.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[0]

  depends_on = [aws_internet_gateway.kumo-internet-gateway]
}

resource "aws_route_table" "kumo-route-table" {
  vpc_id = aws_vpc.kumo-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.kumo-internet-gateway.id
  }

  depends_on = [aws_vpc.kumo-vpc, aws_internet_gateway.kumo-internet-gateway]
}

resource "aws_route_table_association" "kumo-route-table-association" {
  route_table_id = aws_route_table.kumo-route-table.id
  subnet_id      = aws_subnet.kumo-subnet.id

  depends_on = [aws_route_table.kumo-route-table, aws_subnet.kumo-subnet]
}

resource "aws_security_group" "kumo-security-group" {
  name       = "kumo-security-group"
  vpc_id     = aws_vpc.kumo-vpc.id
  depends_on = [aws_vpc.kumo-vpc, aws_route_table_association.kumo-route-table-association]
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

  cidr_ipv4   = local.allowed_ip
  from_port   = 22
  ip_protocol = "tcp"
  to_port     = 22
}

resource "aws_instance" "kumo-ec2-instance" {
  instance_type          = local.AWS_INSTANCE_TYPE
  ami                    = data.aws_ami.kumo-ami.id
  vpc_security_group_ids = [aws_security_group.kumo-security-group.id]
  subnet_id              = aws_subnet.kumo-subnet.id
  availability_zone      = local.first_available_zone

  root_block_device {
    volume_type = local.AWS_EC2_INSTANCE_VOLUME_TYPE
    volume_size = local.AWS_EC2_INSTANCE_VOLUME_SIZE
  }

  depends_on = [aws_vpc.kumo-vpc, aws_subnet.kumo-subnet, aws_security_group.kumo-security-group]
}

output "public_ip" {
  description = "The public IP address of the EC2 instance you can connect to via SSH"
  value       = "Deployment complete, now try 'ssh -i <path_to_your_private_key> <AMI.User>@${aws_instance.kumo-ec2-instance.public_ip}' to connect to your EC2 instance"
}