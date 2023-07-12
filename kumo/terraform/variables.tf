variable "AMI_ID" {
  description = "The AMI ID to use for the EC2 instance"
  type = string
  default = null

  validation {
    condition = length(var.AMI_ID) > 0
    error_message = "AMI_ID must be set"
  }

  validation {
    condition = can(regex("^ami-[0-9a-f]{8,}$", var.AMI_ID))
    error_message = "AMI_ID must be a valid AMI ID"
  }
}

variable "AWS_INSTANCE_TYPE" {
  description = "The EC2 instance type to use"
  type = string
  default = "t2.micro"

  validation {
    condition = length(var.AWS_INSTANCE_TYPE) > 0
    error_message = "AWS_INSTANCE_TYPE must be present"
  }
}

variable "AWS_EC2_INSTANCE_VOLUME_TYPE" {
  description = "The EC2 instance volume type to use"
  type = string
  default = "gp2"

  validation {
    condition = length(var.AWS_EC2_INSTANCE_VOLUME_TYPE) > 0
    error_message = "AWS_EC2_INSTANCE_VOLUME_TYPE must be present"
  }
}

variable "AWS_EC2_INSTANCE_VOLUME_SIZE" {
  description = "The EC2 instance volume size to use"
  type = number
  default = 8

  validation {
    condition = var.AWS_EC2_INSTANCE_VOLUME_SIZE >= 8
    error_message = "AWS_EC2_INSTANCE_VOLUME_SIZE must be greater than or equal to 8"
  }
}