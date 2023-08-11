package iota

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) Name() string {
	switch t {
	case Packer:
		return "packer"
	case Terraform:
		return "terraform"
	default:
		panic("Unknown tool")
	}
}

func (t Tool) VarsName() string {
	switch t {
	case Packer:
		return ".auto.pkrvars.hcl"
	case Terraform:
		return ".auto.tfvars"
	default:
		panic("Unknown tool")
	}
}

func (t Tool) Version() string {
	switch t {
	case Packer:
		return "1.6.5"
	case Terraform:
		return "1.5.3"
	default:
		panic("Unknown tool")
	}
}