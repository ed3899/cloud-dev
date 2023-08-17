variable "AWS_ACCESS_KEY" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS access key ID that Packer will use."

  validation {
    condition     = length(var.AWS_ACCESS_KEY) == 20
    error_message = "The AWS access key should be a 20-character string."
  }

  validation {
    condition     = can(regex("^AKI[0-9A-Z]{17}$", var.AWS_ACCESS_KEY))
    error_message = "The AWS access key value should start with 'AKI' followed by 17 alphanumeric characters."
  }
}

variable "AWS_SECRET_KEY" {
  type        = string
  default     = null
  sensitive   = true
  description = "The AWS secret access key that Packer will use."

  validation {
    condition     = length(var.AWS_SECRET_KEY) == 40
    error_message = "The AWS secret access key should be a 40-character string."
  }

  validation {
    condition     = can(regex("^[0-9a-zA-Z]{40}$", var.AWS_SECRET_KEY, ))
    error_message = "The AWS secret access key should consist only of alphanumeric characters."
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

variable "AWS_USER_IDS" {
  type        = list(string)
  sensitive   = true
  default     = null
  description = "The AWS user ids that can access the instance"

  validation {
    condition     = length(var.AWS_USER_IDS) > 0
    error_message = "The AWS user ids array must not be empty."
  }
}

variable "AWS_AMI_NAME" {
  type        = string
  default     = "cloud_dev"
  description = "The name of the Amazon Machine Image (AMI) that Packer will use."

  validation {
    condition     = length(var.AWS_AMI_NAME) > 0 && length(var.AWS_AMI_NAME) < 129
    error_message = "The AWS AMI name must be between 1 and 128 characters long."
  }
}

variable "AWS_INSTANCE_TYPE" {
  type        = string
  default     = "t2.micro"
  description = "The instance type of the EC2 instance that Packer will create"

  validation {
    condition     = var.AWS_INSTANCE_TYPE != null
    error_message = "Please provide a valid AWS instance type (i.e 't2.micro')."
  }
}

variable "AWS_REGION" {
  type        = string
  default     = "us-west-2"
  description = "The AWS region where Packer will create the resources"

  validation {
    condition = can(regex("^[a-z]+(?:-[a-z]+)*-\\d+$", var.AWS_REGION))

    error_message = "Invalid AWS region format. Must be lowercase letters and dashes only (i.e 'us-east-2')."
  }
}

variable "AWS_EC2_INSTANCE_USERNAME" {
  type        = string
  default     = "dev"
  description = "The username for the EC2 instance you will use to ssh into the machine."

  validation {
    condition     = length(var.AWS_EC2_INSTANCE_USERNAME) > 0
    error_message = "Please provide a username for the EC2 instance."
  }

  validation {
    condition     = length(regexall("^[a-zA-Z_][a-zA-Z0-9_]+$", var.AWS_EC2_INSTANCE_USERNAME)) > 0
    error_message = "The username must contain only alphanumeric characters and '_' (underscore)."
  }
}

variable "AWS_EC2_INSTANCE_USERNAME_PASSWORD" {
  type        = string
  default     = "test123"
  description = "The password for the EC2 instance user."
  sensitive   = true

  validation {
    condition     = length(var.AWS_EC2_INSTANCE_USERNAME_PASSWORD) >= 8 && length(var.AWS_EC2_INSTANCE_USERNAME_PASSWORD) <= 20
    error_message = "The password must be between 8 and 20 characters long without spaces around."
  }
}

variable "AWS_EC2_INSTANCE_USERNAME_HOME" {
  type        = string
  default     = "home"
  description = "The home directory of the EC2 instance user"

  validation {
    condition     = length(var.AWS_EC2_INSTANCE_USERNAME_HOME) > 0
    error_message = "The home directory cannot be an empty string."
  }

  validation {
    condition     = can(regex("^([a-zA-Z0-9]+)$", var.AWS_EC2_INSTANCE_USERNAME_HOME))
    error_message = "The home directory can only contain alphanumeric characters."
  }
}

variable "AWS_EC2_SSH_USERNAME" {
  type        = string
  default     = null
  description = "The SSH username used to initially log into the machine and provision it. This usually changes depending on the public AMI used to build yours"

  validation {
    condition     = length(var.AWS_EC2_SSH_USERNAME) > 0
    error_message = "The SSH username must not be empty."
  }
}

variable "AWS_EC2_AMI_NAME_FILTER" {
  type        = string
  default     = null
  description = "The AMI name that will be used to look for a public AMIs on AWS to build on top of (i.e ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516)"

  validation {
    condition     = length(var.AWS_EC2_AMI_NAME_FILTER) > 0
    error_message = "The AMI name filter must not be empty."
  }
}

variable "AWS_EC2_AMI_ROOT_DEVICE_TYPE" {
  type        = string
  default     = null
  description = "The AMI root device type that will be used to filter out public AMIs on AWS"

  validation {
    condition     = length(var.AWS_EC2_AMI_ROOT_DEVICE_TYPE) > 0
    error_message = "The AMI root device type must not be empty."
  }
}

variable "AWS_EC2_AMI_VIRTUALIZATION_TYPE" {
  type        = string
  default     = null
  description = "The AMI virtualization type that will be used to filter out public AMIs on AWS"

  validation {
    condition     = length(var.AWS_EC2_AMI_VIRTUALIZATION_TYPE) > 0
    error_message = "The AMI virtualization type must not be empty."
  }
}

variable "AWS_EC2_AMI_OWNERS" {
  type        = list(string)
  default     = null
  description = "The AMI virtualization type that will be used to filter out public AMIs on AWS"

  validation {
    condition     = length(var.AWS_EC2_AMI_OWNERS) > 0
    error_message = "The AMI owners array must not be empty."
  }
}

variable "AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL" {
  type        = string
  default     = "/tmp/ansible"
  description = "The directory where ansible files will be uploaded. Packer requires write permissions in this directory."

  validation {
    condition     = can(regex("^/.*", var.AWS_EC2_ANSIBLE_STAGING_DIRECTORY_INTERNAL))
    error_message = "The ansible staging directory must contain an absolute path starting with '/'."
  }
}

variable "AWS_EC2_PUBLIC_DIRECTORY_INTERNAL" {
  type        = string
  default     = "/public"
  description = "The directory where temporary tools are downloaded"

  validation {
    condition     = substr(var.AWS_EC2_PUBLIC_DIRECTORY_INTERNAL, 0, 1) == "/"
    error_message = "The public directory must start with a forward slash (/)."
  }
}

variable "GIT_USERNAME" {
  type        = string
  default     = null
  description = "The git username that will be associated with your commits."

  validation {
    condition     = length(var.GIT_USERNAME) <= 20
    error_message = "Git username must be less than or equal to 20 characters long."
  }
}

variable "GIT_EMAIL" {
  type        = string
  default     = null
  description = "The git email that will be associated with your commits."

  validation {
    condition     = can(regex("^\\S+@\\S+\\.\\S+$", var.GIT_EMAIL))
    error_message = "Git email must be a valid email address."
  }
}

variable "ANSIBLE_TAGS" {
  type        = list(string)
  default     = null
  description = "The ansible tags that will be used to install playbooks."

  validation {
    condition     = length(var.ANSIBLE_TAGS) > 0
    error_message = "The variable must contain at least one tag."
  }
}

variable "GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC" {
  type        = string
  default     = null
  sensitive   = true
  description = "The github personal access token use to setup ssh access on your instance"
}