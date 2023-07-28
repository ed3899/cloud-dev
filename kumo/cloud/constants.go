package cloud

type Tool int

const (
	Packer Tool = iota
	Terraform
)

type Cloud int

const (
	AWS Cloud = iota
)
