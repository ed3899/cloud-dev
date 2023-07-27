package binaries

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

type Cloud int

const (
	AWS Cloud = iota
)

const (
	DEPENDENCIES_DIR_NAME = "deps"
	TEMPLATE_DIR_NAME     = "templates"
	AWS_SUBDIR_NAME       = "aws"
	GENERAL_SUBDIR_NAME   = "general"
	DEFAULT_IP            = "0.0.0.0"
)
