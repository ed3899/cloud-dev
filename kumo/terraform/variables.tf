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
