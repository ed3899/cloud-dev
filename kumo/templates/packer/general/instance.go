package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/templates"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   templates.EnvironmentI
}

func NewTemplate() (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_general_template_failed")
		templatesDirName          = dirs.TEMPLATES_DIR_NAME
		packerDirName             = tool.PACKER_NAME
		generalDirName            = dirs.GENERAL_DIR_NAME
		packerGeneralTemplateName = templates.PACKER_GENERAL_TEMPLATE_NAME

		packerGeneralTemplateInstance  *template.Template
		absPathToPackerGeneralTemplate string
	)

	if absPathToPackerGeneralTemplate, err = filepath.Abs(filepath.Join(templatesDirName, packerDirName, generalDirName, packerGeneralTemplateName)); err != nil {
		err = oopsBuilder.
			With("templatesDirName", templatesDirName).
			With("packerDirName", packerDirName).
			With("generalDirName", generalDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", packerGeneralTemplateName)
		return
	}

	if packerGeneralTemplateInstance, err = template.ParseFiles(absPathToPackerGeneralTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToPackerGeneralTemplate)
		return
	}

	newTemplate = &Template{
		instance:      packerGeneralTemplateInstance,
		parentDirName: packerDirName,
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
