package tool

const (
	PACKER_NAME             = "packer"
	PACKER_VERSION          = "1.9.1"
	PACKER_MANIFEST_NAME    = "manifest.json"
	PACKER_PLUGIN_PATH_NAME = "PACKER_PLUGIN_PATH"

	TERRAFORM_NAME               = "terraform"
	TERRAFORM_VERSION            = "1.5.3"
	TERRAFORM_DEFAULT_ALLOWED_IP = "0.0.0.0"
)

type ToolType int

const (
	Packer ToolType = iota
	Terraform
)
