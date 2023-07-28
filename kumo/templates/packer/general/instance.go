package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type PackerGeneralEnvironment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerGeneralTemplate struct {
	instance    *template.Template
	environment *PackerGeneralEnvironment
}

func NewPackerGeneralTemplate() (packerGeneralTemplate *PackerGeneralTemplate, err error) {
	const (
		PACKER_GENERAL_TEMPLATE_NAME = "GeneralPackerVars.tmpl"
	)

	var (
		oopsBuilder = oops.
				Code("new_packer_general_template_failed")
		packerGeneralTemplateInstance  *template.Template
		absPathToPackerGeneralTemplate string
	)

	if absPathToPackerGeneralTemplate, err = filepath.Abs(filepath.Join(templates.PACKER_TEMPLATES_DIR_NAME, templates.GENERAL_TEMPLATES_DIR_NAME, PACKER_GENERAL_TEMPLATE_NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_TEMPLATES_DIR_NAME", templates.PACKER_TEMPLATES_DIR_NAME).
			With("templates.GENERAL_TEMPLATES_DIR_NAME", templates.GENERAL_TEMPLATES_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", PACKER_GENERAL_TEMPLATE_NAME)
		return
	}

	if packerGeneralTemplateInstance, err = template.ParseFiles(absPathToPackerGeneralTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToPackerGeneralTemplate)
		return
	}

	packerGeneralTemplate = &PackerGeneralTemplate{
		instance: packerGeneralTemplateInstance,
		environment: &PackerGeneralEnvironment{
			GIT_USERNAME:                          viper.GetString("Git.Username"),
			GIT_EMAIL:                             viper.GetString("Git.Email"),
			ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		},
	}

	return
}

func (pgt *PackerGeneralTemplate) GetName() (name string) {
	return pgt.instance.Name()
}

func (pgt *PackerGeneralTemplate) GetInstance() (instance *template.Template) {
	return pgt.instance
}

func (pgt *PackerGeneralTemplate) GetEnvironment() (environment *PackerGeneralEnvironment) {
	return pgt.environment
}
