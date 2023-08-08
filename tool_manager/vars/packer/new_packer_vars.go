package packer

import (
	environment_packer "github.com/ed3899/kumo/tool_manager/environment/packer"
	template_packer "github.com/ed3899/kumo/tool_manager/template/packer"
)

type PackerVars struct {
	Environment *environment_packer.PackerEnvironment
	Template    *template_packer.PackerTemplateCombo
}

func NewPackerVars(
	options ...Option,
) (
	packerVars *PackerVars,
) {
	var (
		option Option
	)

	packerVars = &PackerVars{}

	for _, option = range options {
		option(packerVars)
	}

	return
}

func (pv *PackerVars) OutputVars() (err error) {
	return
}

type Option func(*PackerVars)
