variable "AWS_ACCESS_KEY_ID" {
  description = "The AWS access key ID to use"
  type        = string

  validation {
    condition     = length(var.AWS_ACCESS_KEY_ID) > 0
    error_message = "AWS_ACCESS_KEY_ID must be present"
  }
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "The AWS secret access key to use"
  type        = string

  validation {
    condition     = length(var.AWS_SECRET_ACCESS_KEY) > 0
    error_message = "AWS_SECRET_ACCESS_KEY must be present"
  }
}

variable "AWS_REGION" {
  description = "The AWS region to use"
  type        = string

  validation {
    condition     = length(var.AWS_REGION) > 0
    error_message = "AWS_REGION must be present"
  }
}

variable "AMI_ID" {
  description = "The AMI ID to use for the EC2 instance"
  type        = string

  validation {
    condition     = length(var.AMI_ID) > 0
    error_message = "AMI_ID must be set"
  }

  validation {
    condition     = can(regex("^ami-[0-9a-f]{8,}$", var.AMI_ID))
    error_message = "AMI_ID must be a valid AMI ID"
  }
}

variable "AWS_INSTANCE_TYPE" {
  description = "The EC2 instance type to use"
  type        = string

  validation {
    condition     = length(var.AWS_INSTANCE_TYPE) > 0
    error_message = "AWS_INSTANCE_TYPE must be present"
  }
}

variable "AWS_EC2_INSTANCE_VOLUME_TYPE" {
  description = "The EC2 instance volume type to use"
  type        = string
  default     = "gp2"
}

variable "AWS_EC2_INSTANCE_VOLUME_SIZE" {
  description = "The EC2 instance volume size to use"
  type        = number
  default     = 8

  validation {
    condition     = var.AWS_EC2_INSTANCE_VOLUME_SIZE >= 8
    error_message = "AWS_EC2_INSTANCE_VOLUME_SIZE must be greater than or equal to 8"
  }
}

variable "ALLOWED_IP" {
  description = "The IP address to allow SSH access from"
  type        = string

  validation {
    condition     = length(var.ALLOWED_IP) > 0
    error_message = "ALLOWED_IP must be present"
  }

  validation {
    condition     = can(regex("^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}(?:/[0-9]{1,2})?$", var.ALLOWED_IP))
    error_message = "ALLOWED_IP must be a valid IP address with a CIDR mask"
  }
}

variable "KEY_NAME" {
  description = "The name of the SSH key to create"
  type        = string

  validation {
    condition     = length(var.KEY_NAME) > 0
    error_message = "KEY_NAME must be present"
  }
}

variable "SSH_PORT" {
  description = "The port to use for SSH"
  type        = number

  validation {
    condition     = var.SSH_PORT > 0
    error_message = "SSH_PORT must be greater than 0"
  }
}

variable "IP_FILE_NAME" {
  description = "The file to write the EC2 instance IP address to"
  type        = string

  validation {
    condition     = length(var.IP_FILE_NAME) > 0
    error_message = "IP_FILE_NAME must be present"
  }
}

variable "USERNAME" {
  description = "The username to use for SSH"
  type        = string

  validation {
    condition     = length(var.USERNAME) > 0
    error_message = "USERNAME must be present"
  }
}