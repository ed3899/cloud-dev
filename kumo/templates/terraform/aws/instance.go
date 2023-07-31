package aws

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/templates"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	absPath       string
	parentDirName string
	environment   templates.EnvironmentI
}

func NewTemplate(packerManifest templates.PackerManifestI) (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_template_failed")
		terraformDirName         = tool.TERRAFORM_NAME
		awsDirName               = cloud.AWS_NAME
		terraformAwsTemplateName = templates.TERRAFORM_AWS_TEMPLATE_NAME

		absPathToTerraformTemplate string
		instance                   *template.Template
	)

	if absPathToTerraformTemplate, err = filepath.Abs(filepath.Join(terraformDirName, awsDirName, terraformAwsTemplateName)); err != nil {
		err = oopsBuilder.
			With("terraformDirName", terraformDirName).
			With("awsDirName", awsDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", terraformAwsTemplateName)
		return
	}

	if instance, err = template.ParseFiles(absPathToTerraformTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTerraformTemplate)
		return
	}

	newTemplate = &Template{
		instance:      instance,
		absPath:       absPathToTerraformTemplate,
		parentDirName: terraformDirName,
		environment: &Environment{
			AWS_REGION:                   viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
			AMI_ID:                       packerManifest.GetLastBuiltAmiId(),
		},
	}

	return
}

func (t *Template) GetAbsPath() (absPath string) {
	return t.absPath
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
