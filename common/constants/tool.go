package constants

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) String() (s string) {
	switch t {
	case Packer:
		s = "packer"
	case Terraform:
		s = "terraform"
	}

	return
}

func (t Tool) Version() (v string) {
	switch t {
	case Packer:
		v = "1.6.5"
	case Terraform:
		v = "1.5.3"
	}

	return
}

func (t Tool) Vars() (v string) {
	switch t {
	case Packer:
		v = ".auto.pkrvars.hcl"
	case Terraform:
		v = ".auto.tfvars"
	}

	return
}

func (t Tool) PluginPathEnvironmentVariable(ppev string) {
	switch t {
	case Packer:
		ppev = "PACKER_PLUGIN_PATH"
	case Terraform:
		ppev = "TF_PLUGIN_CACHE_DIR"
	}

	return
}
