package cloud

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

type Kind int

const (
	AWS Kind = iota
)
