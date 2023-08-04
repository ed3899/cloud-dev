package tool

import (
	constants "github.com/ed3899/kumo/0_constants"
	cloud "github.com/ed3899/kumo/1_cloud"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type TerraformTemplate struct {
	General TerraformGeneralEnvironment
	Cloud   CloudEnvironmentI
}

type ToolWithTerraformTemplate struct {
	ToolWithRunDir
	TerraformTemplate
}

func NewToolWithTerraformTemplate(
	toolWithRunDir ToolWithRunDir,
	cloud cloud.Cloud,
	kumoExecutableAbsPath string,
	pickedIp string,
	pickedAmiId string,
) (
	toolWithTerraformTemplate ToolWithTerraformTemplate,
	err error,
) {

	var (
		oopsBuilder = oops.
				Code("NewToolWithTerraformTemplate").
				With("cloud", cloud.Name).
				With("kumoExecutableAbsPath", kumoExecutableAbsPath).
				With("toolWithRunDir", toolWithRunDir)

		terraformGeneralEnvironment TerraformGeneralEnvironment
	)

	terraformGeneralEnvironment = TerraformGeneralEnvironment{
		Required: TerraformGeneralRequired{
			ALLOWED_IP: pickedIp,
		},
	}

	switch cloud.Name {
	case constants.AWS:
		toolWithTerraformTemplate = ToolWithTerraformTemplate{
			ToolWithRunDir: toolWithRunDir,
			TerraformTemplate: TerraformTemplate{
				General: terraformGeneralEnvironment,
				Cloud: TerraformAwsEnvironment{
					Required: TerraformAwsRequired{
						AWS_REGION:        viper.GetString("AWS.Region"),
						AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
						AMI_ID:            pickedAmiId,
						KEY_NAME:          constants.KEY_NAME,
						SSH_PORT:          constants.SSH_PORT,
						IP_FILE_NAME:      constants.IP_FILE_NAME,
						USERNAME:          viper.GetString("AMI.User"),
					},
				},
			},
		}

	default:
		err = oopsBuilder.
			Wrapf(err, "Cloud %s is not supported", cloud.Name)
		return
	}

	return
}
