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