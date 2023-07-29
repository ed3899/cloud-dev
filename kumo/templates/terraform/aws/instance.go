package aws

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Environment struct {
	AWS_REGION                   string
	AWS_INSTANCE_TYPE            string
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
	AMI_ID                       string
}

func (e *Environment) IsEnvironment() (isEnvironment bool) {
	return true
}

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   *Environment
}

func NewTemplate(amiIdToBeUsed string) (newTemplate *Template, err error) {
	const (
		NAME       = "AWS_TerraformVars.tmpl"
		DEFAULT_IP = "0.0.0.0"
	)

	var (
		oopsBuilder = oops.
				Code("new_template_failed")

		absPath  string
		instance *template.Template
	)

	if absPath, err = filepath.Abs(filepath.Join(templates.TERRAFORM_DIR_NAME, templates.AWS_DIR_NAME, NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_DIR_NAME", templates.PACKER_DIR_NAME).
			With("templates.AWS_DIR_NAME", templates.AWS_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", NAME)
		return
	}

	if instance, err = template.ParseFiles(absPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPath)
		return
	}

	newTemplate = &Template{
		instance:      instance,
		parentDirName: templates.TERRAFORM_DIR_NAME,
		environment: &Environment{
			AWS_REGION:                   viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
			AMI_ID:                       amiIdToBeUsed,
		},
	}

	return
}
