package tool

import (
	constants "github.com/ed3899/kumo/0_constants"
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type PackerTemplate struct {
	General PackerGeneralEnvironment
	Cloud   CloudEnvironmentI
}

type ToolWithPackerTemplate struct {
	ToolWithRunDir
	PackerTemplate
}

func NewToolWithPackerTemplate(
	toolWithRunDir ToolWithRunDir,
	cloud cloud.Cloud,
	kumoExecutableAbsPath string,
) (toolWithPackerTemplate ToolWithPackerTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("NewToolWithPackerTemplate").
				With("cloud", cloud).
				With("kumoExecutableAbsPath", kumoExecutableAbsPath)

		packerGeneralEnvironment PackerGeneralEnvironment
	)

	packerGeneralEnvironment = PackerGeneralEnvironment{
		Required: PackerGeneralRequired{
			GIT_USERNAME: viper.GetString("Git.Username"),
			GIT_EMAIL:    viper.GetString("Git.Email"),
			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
		},
		Optional: PackerGeneralOptional{
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		},
	}

	switch cloud.Name {
	case constants.AWS:
		toolWithPackerTemplate = ToolWithPackerTemplate{
			ToolWithRunDir: toolWithRunDir,
			PackerTemplate: PackerTemplate{
				General: packerGeneralEnvironment,
				Cloud: PackerAwsEnvironment{
					Required: PackerAwsRequired{
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
			},
		}
	default:
		err = oopsBuilder.
			Errorf("Unknown cloud: %s", cloud.Name)
		return
	}

	return
}
