package environment

import (
	"github.com/ed3899/kumo/constants"
	"github.com/spf13/viper"
)

func NewTerraformAwsEnvironment(pickedAmiId string) (terraformAwsEnvironment CloudEnvironmentI) {
	terraformAwsEnvironment = TerraformAwsEnvironment{
		Required: TerraformAwsRequired{
			AWS_REGION:        viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
			AMI_ID:            pickedAmiId,
			KEY_NAME:          constants.KEY_NAME,
			SSH_PORT:          constants.SSH_PORT,
			IP_FILE_NAME:      constants.IP_FILE_NAME,
			USERNAME:          viper.GetString("AMI.User"),
		},
		Optional: TerraformAwsOptional{
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
		},
	}
	return
}

type TerraformAwsRequired struct {
	AWS_REGION        string
	AWS_INSTANCE_TYPE string
	AMI_ID            string
	KEY_NAME          string
	SSH_PORT          int
	IP_FILE_NAME      string
	USERNAME          string
}

type TerraformAwsOptional struct {
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}

type TerraformAwsEnvironment struct {
	Required TerraformAwsRequired
	Optional TerraformAwsOptional
}

func (tae TerraformAwsEnvironment) IsEnvironment() bool {
	return true
}

func (tae TerraformAwsEnvironment) IsCloudEnvironment() bool {
	return true
}
