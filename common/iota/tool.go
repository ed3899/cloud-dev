package iota

import (
	"log"

	"github.com/samber/oops"
)

type Tool int

const (
	Packer Tool = iota
	Terraform
)

func (t Tool) Name() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Name")

	switch t {
	case Packer:
		return "packer"

	case Terraform:
		return "terraform"

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		log.Fatalf("%+v", err)

		return ""
	}
}

func (t Tool) VarsName() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("VarsName")

	switch t {
	case Packer:
		return ".auto.pkrvars.hcl"

	case Terraform:
		return ".auto.tfvars"

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		log.Fatalf("%+v", err)

		return ""
	}
}

func (t Tool) Version() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Version")

	switch t {
	case Packer:
		return "1.6.5"

	case Terraform:
		return "1.5.3"

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		log.Fatalf("%+v", err)

		return ""
	}
}

func (t Tool) PluginPathEnvironmentVariable() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("PluginPathEnvironmentVariable")

	switch t {
	case Packer:
		return "PACKER_PLUGIN_PATH"

	case Terraform:
		return "TF_PLUGIN_CACHE_DIR"

	default:
		err := oopsBuilder.
			Errorf("unknown tool: %#v", t)

		log.Fatalf("%+v", err)

		return ""
	}
}
