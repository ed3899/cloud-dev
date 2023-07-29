package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   templates.EnvironmentI
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

	if absPathToPackerGeneralTemplate, err = filepath.Abs(filepath.Join(dirs.PACKER_DIR_NAME, dirs.GENERAL_DIR_NAME, PACKER_GENERAL_TEMPLATE_NAME)); err != nil {
		err = oopsBuilder.
			With("dirs.PACKER_DIR_NAME", dirs.PACKER_DIR_NAME).
			With("dirs.GENERAL_DIR_NAME", dirs.GENERAL_DIR_NAME).
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
		parentDirName: dirs.PACKER_DIR_NAME,
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

func (t *Template) GetEnvironment() (environment templates.EnvironmentI) {
	return t.environment
}
