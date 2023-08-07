package aws

import (
	utils_environment "github.com/ed3899/kumo/utils/environment"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewEnvironment() (environment Environment) {
	environment = Environment{
		Required: Required{
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

type Required struct {
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

type Environment struct {
	Required Required
}

type NewEnvironmentF func() Environment

func (e Environment) IsCloudEnvironment() (isCloudEnvironment bool) {
	isCloudEnvironment = true
	return
}

type Option func(Environment) Environment

func NewEnv(
	requiredFieldsAreNotFilled utils_environment.IsStructNotCompletelyFilledF,
	opts ...Option,
) (environment Environment, err error) {

	var (
		oopsBuilder = oops.
				Code("NewEnv").
				With("opts", opts)

		notFilled bool
		missingField string
	)

	environment = Environment{}
	for _, o := range opts {
		environment = o(environment)
	}

	notFilled, missingField = requiredFieldsAreNotFilled(environment.Required)
	if notFilled {
		err = oopsBuilder.
			With("environment.Required", environment.Required).
			Errorf("Required field '%s' is not filled", missingField)
		return
	}

	return
}
