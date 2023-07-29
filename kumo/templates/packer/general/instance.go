package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Environment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

func (pge *Environment) IsEnvironment() (isEnvironment bool) {
	return true
}

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   *Environment
}

func NewTemplate() (newTemplate *Template, err error) {
	const (
		PACKER_GENERAL_TEMPLATE_NAME = "GeneralPackerVars.tmpl"
	)

	var (
		oopsBuilder = oops.
				Code("new_packer_general_template_failed")
		packerGeneralTemplateInstance  *template.Template
		absPathToPackerGeneralTemplate string
	)

	if absPathToPackerGeneralTemplate, err = filepath.Abs(filepath.Join(templates.PACKER_DIR_NAME, templates.GENERAL_DIR_NAME, PACKER_GENERAL_TEMPLATE_NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_DIR_NAME", templates.PACKER_DIR_NAME).
			With("templates.GENERAL_TEMPLATES_DIR_NAME", templates.GENERAL_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", PACKER_GENERAL_TEMPLATE_NAME)
		return
	}

	if packerGeneralTemplateInstance, err = template.ParseFiles(absPathToPackerGeneralTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToPackerGeneralTemplate)
		return
	}

	newTemplate = &Template{
		instance:      packerGeneralTemplateInstance,
		parentDirName: templates.PACKER_DIR_NAME,
		environment: &Environment{
			GIT_USERNAME:                          viper.GetString("Git.Username"),
			GIT_EMAIL:                             viper.GetString("Git.Email"),
			ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		},
	}

	return
}

func (t *Template) GetParentDirName() (dir string) {
	return t.parentDirName
}

func (t *Template) GetName() (name string) {
	return t.instance.Name()
}

func (t *Template) GetInstance() (instance *template.Template) {
	return t.instance
}

func (t *Template) GetEnvironment() (environment *Environment) {
	return t.environment
}
