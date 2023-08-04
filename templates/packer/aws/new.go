package aws

import (
	"path/filepath"
	"text/template"

	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/templates"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	absPath       string
	parentDirName string
	environment   templates.EnvironmentI
}

func New() (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_template_failed").
				With("templates.PACKER_AWS_TEMPLATE_NAME", templates.PACKER_AWS_TEMPLATE_NAME)
		templatesDirName      = dirs.TEMPLATES_DIR_NAME
		packerDirName         = common_tool_constants.PACKER_NAME
		awsDirName            = common_cloud_constants.AWS_NAME
		packerAwsTemplateName = templates.PACKER_AWS_TEMPLATE_NAME

		instance          *template.Template
		absPathToTemplate string
	)

	if absPathToTemplate, err = filepath.Abs(filepath.Join(templatesDirName, packerDirName, awsDirName, packerAwsTemplateName)); err != nil {
		err = oopsBuilder.
			With("templatesDirName", templatesDirName).
			With("packerDirName", packerDirName).
			With("awsDirName", awsDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", packerAwsTemplateName)
		return
	}

	if instance, err = template.ParseFiles(absPathToTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTemplate)
		return
	}

	newTemplate = &Template{
		instance:      instance,
		absPath:       absPathToTemplate,
		parentDirName: packerDirName,
		environment: &Environment{
			Required: &Required{
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

func (t *Template) Environment() (environment templates.EnvironmentI) {
	return t.environment
}
