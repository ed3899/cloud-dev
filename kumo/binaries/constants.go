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
	AWS_SUBDIR_NAME = "aws"
)
