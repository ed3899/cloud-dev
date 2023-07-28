package tool

type Tool int

const (
	Packer Tool = iota
	Terraform
)

const (
	PACKER_RUN_DIR_NAME = "packer"
	TERRAFORM_RUN_DIR_NAME = "terraform"
)