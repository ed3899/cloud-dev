package iota

import "github.com/samber/oops"

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) Name() string {
	var choice string

	oops.
		In("iota").
		Tags("tool").
		Code("Name").
		Recoverf(func() {
			switch t {
			case Packer:
				choice = "packer"
			case Terraform:
				choice = "terraform"
			default:
				panic(t)
			}
		}, "Unknown tool")

	return choice
}

func (t Tool) VarsName() string {
	var choice string

	oops.
		In("iota").
		Tags("tool").
		Code("VarsName").
		Recoverf(func() {
			switch t {
			case Packer:
				choice = ".auto.pkrvars.hcl"
			case Terraform:
				choice = ".auto.tfvars"
			default:
				panic(t)
			}
		}, "Unknown tool")

	return choice
}

func (t Tool) Version() string {
	var choice string

	oops.
		In("iota").
		Tags("tool").
		Code("Version").
		Recoverf(func() {
			switch t {
			case Packer:
				choice = "1.6.5"
			case Terraform:
				choice = "1.5.3"
			default:
				panic(t)
			}
		}, "Unknown tool")

	return choice
}

func (t Tool) PluginPathEnvironmentVariable() string {
	var choice string

	oops.
		In("iota").
		Tags("tool").
		Code("PluginPathEnvironmentVariable").
		Recoverf(func() {
			switch t {
			case Packer:
				choice = "PACKER_PLUGIN_PATH"
			case Terraform:
				choice = "TF_PLUGIN_CACHE_DIR"
			default:
				panic(t)
			}
		}, "Unknown tool")

	return choice
}
