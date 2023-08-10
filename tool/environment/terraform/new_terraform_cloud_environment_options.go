package environment

import (
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool_manager/environment/terraform/aws"
	"github.com/ed3899/kumo/utils/packer_manifest"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewTerraformCloudEnvironmentOptions(
	cloud cloud.Cloud,
	kumoExecAbsPath string,
	selectedAmiIdFromConfig string,
	getLastAmiIdFromPackerManifest packer_manifest.GetLastBuiltAmiIdFromPackerManifestF,
	pickAmiId packer_manifest.PickAmiIdF,
) (
	terraformCloudEnvironmentOptions *TerraformCloudEnvironmentOptions,
	err error,
) {

	var (
		oopsBuilder = oops.
				Code("NewTerraformCloudEnvironmentOptions").
				With("cloud", cloud.Name).
				With("kumoExecAbsPath", kumoExecAbsPath).
				With("getLastAmiIdFromPackerManifest", getLastAmiIdFromPackerManifest).
				With("pickAmiId", pickAmiId)

		packerManifestAbsPath string
		lastBuiltAmiId        string
		pickedAmiId           string
	)

	switch cloud.Kind {
	case constants.Aws:
		packerManifestAbsPath = filepath.Join(kumoExecAbsPath, constants.PACKER, constants.AWS, constants.PACKER_MANIFEST)

		if lastBuiltAmiId, err = getLastAmiIdFromPackerManifest(packerManifestAbsPath); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Failed to get last built AMI ID from packer manifest '%s'", packerManifestAbsPath)
			return
		}

		if pickedAmiId, err = pickAmiId(lastBuiltAmiId, selectedAmiIdFromConfig); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Failed to pick AMI ID from last built AMI ID '%s' and user selected AMI ID '%s'", lastBuiltAmiId, selectedAmiIdFromConfig)
			return
		}

		terraformCloudEnvironmentOptions = &TerraformCloudEnvironmentOptions{
			Aws: &aws.TerraformAwsEnvironment{
				Required: &aws.TerraformAwsRequired{
					AWS_REGION:        viper.GetString("AWS.Region"),
					AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
					AMI_ID:            pickedAmiId,
					KEY_NAME:          constants.KEY_NAME,
					SSH_PORT:          constants.SSH_PORT,
					IP_FILE_NAME:      constants.IP_FILE_NAME,
					USERNAME:          viper.GetString("AMI.User"),
				},
				Optional: &aws.TerraformAwsOptional{
					AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
					AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
				},
			},
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' is not supported", cloud.Name)
		return
	}

	return
}

type TerraformCloudEnvironmentOptions struct {
	Aws *aws.TerraformAwsEnvironment
}
