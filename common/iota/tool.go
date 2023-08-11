package iota

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) String() string {
	switch t {
	case Packer:
		return "packer"
	case Terraform:
		return "terraform"
	default:
		panic("Unknown tool")
	}
}
