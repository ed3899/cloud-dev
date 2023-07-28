package aws

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/templates"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Environment struct {
	AWS_ACCESS_KEY                     string
	AWS_SECRET_KEY                     string
	AWS_IAM_PROFILE                    string
	AWS_USER_IDS                       []string
	AWS_AMI_NAME                       string
	AWS_INSTANCE_TYPE                  string
	AWS_REGION                         string
	AWS_EC2_AMI_NAME_FILTER            string
	AWS_EC2_AMI_ROOT_DEVICE_TYPE       string
	AWS_EC2_AMI_VIRTUALIZATION_TYPE    string
	AWS_EC2_AMI_OWNERS                 []string
	AWS_EC2_SSH_USERNAME               string
	AWS_EC2_INSTANCE_USERNAME          string
	AWS_EC2_INSTANCE_USERNAME_HOME     string
	AWS_EC2_INSTANCE_USERNAME_PASSWORD string
}

func (pae *Environment) IsPackerCloudEnvironment() (isPackerCloudEnvironment bool) {
	return true
}

type PackerAwsTemplate struct {
	instance    *template.Template
	environment *Environment
}

func NewPackerAwsTemplate() (packerAwsTemplate *PackerAwsTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_aws_template_failed")
		packerAwsTemplateInstance  *template.Template
		absPathToPackerAwsTemplate string
	)

	if absPathToPackerAwsTemplate, err = filepath.Abs(filepath.Join(templates.PACKER_TEMPLATES_DIR_NAME, templates.AWS_TEMPLATES_DIR_NAME, PACKER_AWS_TEMPLATE_NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_TEMPLATES_DIR_NAME", templates.PACKER_TEMPLATES_DIR_NAME).
			With("templates.AWS_TEMPLATES_DIR_NAME", templates.AWS_TEMPLATES_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", PACKER_AWS_TEMPLATE_NAME)
		return
	}

	if packerAwsTemplateInstance, err = template.ParseFiles(absPathToPackerAwsTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToPackerAwsTemplate)
		return
	}

	packerAwsTemplate = &PackerAwsTemplate{
		instance: packerAwsTemplateInstance,
		environment: &Environment{
			AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
			AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
			AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
			AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
			AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
			AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
			AWS_REGION:                         viper.GetString("AWS.Region"),
			AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
			AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
			AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
			AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
			AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
			AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
			AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
			AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
		},
	}

	return
}

func (pat *PackerAwsTemplate) GetName() (name string) {
	return pat.instance.Name()
}

func (pat *PackerAwsTemplate) GetInstance() (instance *template.Template) {
	return pat.instance
}

func (pat *PackerAwsTemplate) GetEnvironment() (environment *Environment) {
	return pat.environment
}
