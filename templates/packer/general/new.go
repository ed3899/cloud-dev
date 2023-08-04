package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	common_templates_constants "github.com/ed3899/kumo/common/templates/constants"
	common_templates_interfaces "github.com/ed3899/kumo/common/templates/interfaces"
	common_templates_structs "github.com/ed3899/kumo/common/templates/structs"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	absPath       string
	parentDirName string
	environment   common_templates_interfaces.Environment
}

func New() (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_general_template_failed")
		templatesDirName          = dirs.TEMPLATES_DIR_NAME
		packerDirName             = common_tool_constants.PACKER_NAME
		generalDirName            = dirs.GENERAL_DIR_NAME
		packerGeneralTemplateName = common_templates_constants.PACKER_GENERAL_TEMPLATE_NAME

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
		absPath:       absPathToPackerGeneralTemplate,
		parentDirName: packerDirName,
		environment: &common_templates_structs.PackerGeneralEnvironment{
			Required: &common_templates_structs.PackerGeneralRequired{
				GIT_USERNAME: viper.GetString("Git.Username"),
				GIT_EMAIL:    viper.GetString("Git.Email"),
				ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools")},
			Optional: &common_templates_structs.PackerGeneralOptional{
				GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
			},
		},
	}

	return
}

func (t *Template) AbsPath() (absPath string) {
	return t.absPath
}

func (t *Template) ParentDirName() (dir string) {
	return t.parentDirName
}

func (t *Template) Instance() (instance *template.Template) {
	return t.instance
}

func (t *Template) Environment() (environment common_templates_interfaces.Environment) {
	return t.environment
}
