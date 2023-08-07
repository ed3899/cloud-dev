package aws

import (
	"path/filepath"

	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/packer_manifest"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewEnvironment(
	getLastAmiIdFromPackerManifest packer_manifest.GetLastBuiltAmiIdFromPackerManifestF,
	pickAmiId packer_manifest.PickAmiIdF,
	kumoExecAbsPath string,
) (environment Environment) {
	var (
		oopsBuilder = oops.
				Code("NewEnvironment").
				With("kumoExecAbsPath", kumoExecAbsPath).
				With("getLastAmiIdFromPackerManifest", getLastAmiIdFromPackerManifest).
				With("pickAmiId", pickAmiId)
		packerManifestAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.PACKER,
			constants.AWS,
			constants.PACKER_MANIFEST,
		)

		lastAmiId string
	)

	



	environment = Environment{
		Required: Required{
			AWS_REGION:        viper.GetString("AWS.Region"),
			AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
			AMI_ID:            pickedAmiId,
			KEY_NAME:          constants.KEY_NAME,
			SSH_PORT:          constants.SSH_PORT,
			IP_FILE_NAME:      constants.IP_FILE_NAME,
			USERNAME:          viper.GetString("AMI.User"),
		},
		Optional: Optional{
			AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
			AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
		},
	}

	return
}

type Required struct {
	AWS_REGION        string
	AWS_INSTANCE_TYPE string
	AMI_ID            string
	KEY_NAME          string
	SSH_PORT          int
	IP_FILE_NAME      string
	USERNAME          string
}

type Optional struct {
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}

type Environment struct {
	Required Required
	Optional Optional
}

func (tae Environment) IsCloudEnvironment() bool {
	return true
}
