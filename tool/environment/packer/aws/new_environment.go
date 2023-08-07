package aws

import (
	utils_environment "github.com/ed3899/kumo/utils/environment"
	"github.com/samber/oops"
)

// func NewAwsEnvironment() (environment Environment) {
// 	environment = Environment{
// 		Required: Required{
// 			AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
// 			AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
// 			AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
// 			AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
// 			AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
// 			AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
// 			AWS_REGION:                         viper.GetString("AWS.Region"),
// 			AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
// 			AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
// 			AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
// 			AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
// 			AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
// 			AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
// 			AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
// 			AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
// 		},
// 	}

// 	return
// }

type NewPackerAwsEnvironmentF func() (*PackerAwsEnvironment, error)

func NewPackerAwsEnvironment(
	requiredFieldsAreNotFilled utils_environment.IsStructNotCompletelyFilledF,
	opts ...Option,
) (environment *PackerAwsEnvironment, err error) {

	var (
		oopsBuilder = oops.
				Code("NewEnv").
				With("opts", opts)

		notFilled    bool
		missingField string
	)

	environment = &PackerAwsEnvironment{}
	for _, o := range opts {
		o(environment)
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

func WithAwsAccessKey(awsAccessKey string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_ACCESS_KEY = awsAccessKey
	}

	return
}

func WithAwsSecretKey(awsSecretKey string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_SECRET_KEY = awsSecretKey
	}

	return
}

func WithAwsIamProfile(awsIAmProfile string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_IAM_PROFILE = awsIAmProfile
	}

	return
}

func WithAwsUserIds(awsUserIds []string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_USER_IDS = awsUserIds
	}

	return
}

func WithAwsAmiName(awsAmiName string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_AMI_NAME = awsAmiName
	}

	return
}

func WithAwsInstanceType(awsInstanceType string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_INSTANCE_TYPE = awsInstanceType
	}

	return
}

func WithAwsRegion(awsRegion string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_REGION = awsRegion
	}

	return
}

func WithAwsEc2AmiNameFilter(awsEc2AmiNameFilter string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_AMI_NAME_FILTER = awsEc2AmiNameFilter
	}

	return
}

func WithAwsEc2AmiRootDeviceType(awsEc2AmiRootDeviceType string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_AMI_ROOT_DEVICE_TYPE = awsEc2AmiRootDeviceType
	}

	return
}

func WithAwsEc2AmiVirtualizationType(awsEc2AmiVirtualizationType string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_AMI_VIRTUALIZATION_TYPE = awsEc2AmiVirtualizationType
	}

	return
}

func WithAwsEc2AmiOwners(awsEc2AmiOwners []string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_AMI_OWNERS = awsEc2AmiOwners
	}

	return
}

func WithAwsEc2SshUsername(awsEc2SshUsername string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_SSH_USERNAME = awsEc2SshUsername
	}

	return
}

func WithAwsEc2InstanceUsername(awsEc2InstanceUsername string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_INSTANCE_USERNAME = awsEc2InstanceUsername
	}

	return
}

func WithAwsEc2InstanceUsernameHome(awsEc2InstanceUsernameHome string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_INSTANCE_USERNAME_HOME = awsEc2InstanceUsernameHome
	}

	return
}

func WithAwsEc2InstanceUsernamePassword(awsEc2InstanceUsernamePassword string) (option Option) {
	option = func(environment *PackerAwsEnvironment) {
		environment.Required.AWS_EC2_INSTANCE_USERNAME_PASSWORD = awsEc2InstanceUsernamePassword
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

type PackerAwsEnvironment struct {
	Required Required
}

func (e *PackerAwsEnvironment) IsCloudEnvironment() (isCloudEnvironment bool) {
	isCloudEnvironment = true
	return
}

type Option func(*PackerAwsEnvironment)
