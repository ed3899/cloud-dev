package iota

import "github.com/samber/oops"

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) Name() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Name")

	switch t {
	case Packer:
		return "packer", nil

	case Terraform:
		return "terraform", nil

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		return "", err
	}
}

func (t Tool) VarsName() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("VarsName")

	switch t {
	case Packer:
		return ".auto.pkrvars.hcl", nil

	case Terraform:
		return ".auto.tfvars", nil
		
	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		return "", err
	}
}

func (t Tool) Version() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Version")

	switch t {
	case Packer:
		return "1.6.5", nil

	case Terraform:
		return "1.5.3", nil

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		return "", err
	}
}

func (t Tool) PluginPathEnvironmentVariable() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("PluginPathEnvironmentVariable")

	switch t {
	case Packer:
		return "PACKER_PLUGIN_PATH", nil

	case Terraform:
		return "TF_PLUGIN_CACHE_DIR", nil

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		return "", err
	}
}
