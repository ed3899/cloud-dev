package environment

import (
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/utils/packer_manifest"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewTerraformAwsEnvironment(
	pathToPackerManifest string,
	cloud iota.Cloud,
) (*TerraformAwsEnvironment, error) {
	oopBuilder := oops.
		Code("NewTerraformAwsEnvironment")

	amiId, err := packer_manifest.GetLastBuiltAmiIdFromPackerManifest(pathToPackerManifest)
	if err != nil {
		return nil, oopBuilder.
			Wrapf(err, "failed to get last built ami id from packer manifest")
	}

	pickedAmiId, err := packer_manifest.PickAmiId(amiId, viper.GetString("Up.AmiId"))
	if err != nil {
		return nil, oopBuilder.
			Wrapf(err, "failed to pick ami id")
	}

	return &TerraformAwsEnvironment{
		Required: &TerraformAwsRequired{
			AWS_ACCESS_KEY_ID:     viper.GetString("AWS.AccessKeyId"),
			AWS_SECRET_ACCESS_KEY: viper.GetString("AWS.SecretAccessKey"),
			AWS_REGION:            viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE:     viper.GetString("AWS.EC2.Instance.Type"),
			AMI_ID:                pickedAmiId,
			KEY_NAME:              constants.KEY_NAME,
			SSH_PORT:              constants.SSH_PORT,
			IP_FILE_NAME:          constants.IP_FILE_NAME,
			USERNAME:              viper.GetString("AMI.User"),
		},
		Optional: &TerraformAwsOptional{
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
		},
	}, nil
}

type TerraformAwsEnvironment struct {
	Required *TerraformAwsRequired
	Optional *TerraformAwsOptional
}
type TerraformAwsRequired struct {
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	AWS_INSTANCE_TYPE     string
	AMI_ID                string
	KEY_NAME              string
	SSH_PORT              int
	IP_FILE_NAME          string
	USERNAME              string
}

type TerraformAwsOptional struct {
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}
