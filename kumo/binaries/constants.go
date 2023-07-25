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
)

const (
	AWS_SUBDIR_NAME = "aws"
)
