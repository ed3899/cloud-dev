package tool

const (
	PACKER_RUN_DIR_NAME    = "packer"
	PACKER_NAME = "packer"
	
	TERRAFORM_RUN_DIR_NAME = "terraform"
	TERRAFORM_NAME = "terraform"
)

type ToolType int

const (
	Packer ToolType = iota
	Terraform
)
