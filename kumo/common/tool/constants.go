package tool

const (
	PACKER_RUN_DIR_NAME = "packer"
	TERRAFORM_RUN_DIR_NAME = "terraform"
)

type Tool int

const (
	Packer Tool = iota
	Terraform
)