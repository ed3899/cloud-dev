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

func (e *Environment) IsEnvironment() (isEnvironment bool) {
	return true
}

type Template struct {
	instance      *template.Template
	parentDirName string
	environment   *Environment
}

func NewTemplate() (newTemplate *Template, err error) {
	const (
		NAME = "AWS_PackerVars.tmpl"
	)

	var (
		oopsBuilder = oops.
				Code("new_template_failed").
				With("TEMPLATE_NAME", NAME)

		instance *template.Template
		absPath  string
	)

	if absPath, err = filepath.Abs(filepath.Join(templates.PACKER_DIR_NAME, templates.AWS_DIR_NAME, NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_DIR_NAME", templates.PACKER_DIR_NAME).
			With("templates.AWS_TEMPLATES_DIR_NAME", templates.AWS_DIR_NAME).
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
		parentDirName: templates.PACKER_DIR_NAME,
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
