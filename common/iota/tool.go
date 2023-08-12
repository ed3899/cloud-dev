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

func (t Tool) Iota() Tool {
	return t
}

func (t Tool) Name() (n string) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Name")

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch t {
	case Packer:
		return "packer"

	case Terraform:
		return "terraform"

	default:
		panic(t)
	}
}

func (t Tool) VarsName() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("VarsName")

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch t {
	case Packer:
		return ".auto.pkrvars.hcl"

	case Terraform:
		return ".auto.tfvars"

	default:
		panic(t)
	}
}

func (t Tool) Version() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("Version")

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch t {
	case Packer:
		return "1.6.5"

	case Terraform:
		return "1.5.3"

	default:
		panic(t)
	}
}

func (t Tool) PluginPathEnvironmentVariable() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Tool").
		Code("PluginPathEnvironmentVariable")

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch t {
	case Packer:
		return "PACKER_PLUGIN_PATH"

	case Terraform:
		return "TF_PLUGIN_CACHE_DIR"

	default:
		panic(t)
	}
}
