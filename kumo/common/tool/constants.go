package tool

const (
	PACKER_RUN_DIR_NAME    = "packer"
	TERRAFORM_RUN_DIR_NAME = "terraform"
)

type Type int

const (
	Packer Type = iota
	Terraform
)
