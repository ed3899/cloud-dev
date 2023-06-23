variable "AWS_ACCESS_KEY" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS access key ID that Packer will use."

  validation {
    condition     = length(var.AWS_ACCESS_KEY) == 20 && can(regex("^AKI[0-9A-Z]{17}$", var.AWS_ACCESS_KEY))
    error_message = "The AWS access key value is invalid. It should be a 20-character string starting with 'AKI' followed by 17 alphanumeric characters."
  }
}

variable "AWS_SECRET_KEY" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS secret access key that Packer will use."

  validation {
    condition     = length(var.AWS_SECRET_KEY) == 40 && can(regex("^[0-9a-zA-Z]{40}$", var.AWS_SECRET_KEY, ))
    error_message = "The AWS secret access key value is invalid. It should be a 40-character string consisting of alphanumeric characters."
  }
}

variable "AWS_IAM_PROFILE" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS IAM profile name that Packer will use."

  validation {
    condition     = can(regex("^([a-zA-Z0-9-_+=,.@]){1,64}$", var.AWS_IAM_PROFILE))
    error_message = "The AWS IAM profile value is invalid. It should be a string consisting of alphanumeric characters and certain special characters (-_+=,.@), with a length between 1 and 64 characters."
  }
}

variable "AWS_USER_ID" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS User ID that Packer will use."

  validation {
    condition     = can(regex("[0-9A-Z]{0,64}$", var.AWS_USER_ID))
    error_message = "The AWS User ID value is invalid. It should be a string consisting of alphanumeric characters."
  }
}

variable "AWS_AMI_NAME" {
  type        = string
  default     = null
  sensitive   = true
  description = "The name of the Amazon Machine Image (AMI) that Packer will use."


  validation {
    condition     = length(var.AWS_AMI_NAME) > 0 && length(var.AWS_AMI_NAME) < 129
    error_message = "The AWS AMI name must be between 1 and 128 characters long."
  }
}

variable "AWS_INSTANCE_TYPE" {
  type        = string
  default     = null
  sensitive   = true
  description = "The instance type of the EC2 instance that Packer will create"

  validation {
    condition     = var.AWS_INSTANCE_TYPE != null
    error_message = "Please provide a valid AWS instance type (i.e 't2.micro')."
  }
}

variable "AWS_REGION" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS region where Packer will create the resources"

  validation {
    condition = can(regex("^[a-z]+(?:-[a-z]+)*-\\d+$", var.AWS_REGION))

    error_message = "Invalid AWS region format. Must be lowercase letters and dashes only (i.e 'us-east-2')."
  }
}

variable "AWS_INSTANCE_SSH_USERNAME" {
  type        = string
  default     = null
  sensitive   = true
  description = "The SSH username for the EC2 instance that Packer will create"

  validation {
    condition     = can(regex("^[a-zA-Z_][a-zA-Z0-9_-]{0,30}[a-zA-Z0-9_]$", var.AWS_INSTANCE_SSH_USERNAME))
    error_message = "Invalid AWS instance SSH username. Must start with a letter or underscore and contain only alphanumeric characters, hyphens or underscores."
  }
}