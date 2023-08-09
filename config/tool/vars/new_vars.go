package vars

import (
	"path/filepath"

	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/config/tool/template"
)

func NewVars(
	options ...Option,
) (
	vars *Vars,
) {
	var (
		option Option
	)

	vars = &Vars{}
	for _, option = range options {
		option(vars)
	}

	return
}

func WithAbsPath(
	kumoExecAbsPath string,
	toolConfig tool.ToolConfig,
	cloudConfig cloud.CloudConfig,
) (
	option Option,
) {
	option = func(vars *Vars) {
		vars.AbsPath = filepath.Join(
			kumoExecAbsPath,
			toolConfig.Name(),
			cloudConfig.Name,
		)
	}

	return
}

func WithTemplateFile(
	templateFile *template.TemplateFile,
) (
	option Option,
) {
	option = func(vars *Vars) {
		vars.TemplateFile = templateFile
	}

	return
}

func (v *Vars) Create() (err error) {
	return
}

type Vars struct {
	AbsPath      string
	TemplateFile *template.TemplateFile
}

type Option func(*Vars)
