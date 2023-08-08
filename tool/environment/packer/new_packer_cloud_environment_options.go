package packer

import (
	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool/environment/packer/aws"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewPackerCloudEnvironmentOptions(
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (
	packerCloudEnvironmentOptions *PackerCloudEnvironmentOptions,
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("NewPackerCloudEnvironmentOptions").
			With("cloud", cloud.Name).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	switch cloud.Kind {
	case constants.Aws:
		packerCloudEnvironmentOptions = &PackerCloudEnvironmentOptions{
			Aws: &aws.PackerAwsEnvironment{
				Required: &aws.PackerAwsRequired{
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

	default:
		err = oopsBuilder.
			Errorf("Unsupported cloud kind '%v'", cloud.Kind)
		return
	}

	return
}

type PackerCloudEnvironmentOptions struct {
	Aws *aws.PackerAwsEnvironment
}
