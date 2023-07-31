package tool

const (
	PACKER_RUN_DIR_NAME  = "packer"
	PACKER_NAME          = "packer"
	PACKER_MANIFEST_NAME = "manifest.json"

	TERRAFORM_RUN_DIR_NAME       = "terraform"
	TERRAFORM_NAME               = "terraform"
	TERRAFORM_DEFAULT_ALLOWED_IP = "0.0.0.0"
)

type ToolType int

const (
	Packer ToolType = iota
	Terraform
)
