package templates

import (
	"log"

	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	templates_terraform_general "github.com/ed3899/kumo/templates/terraform/general"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"

	"github.com/samber/lo"
)

func CraftCloudTerraformTfVarsFile(cloud string) (cloudTfVarsFilePath string, err error) {
	packerManifestPath, err := utils.GetPackerManifestPathTo(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while getting packer manifest path for cloud '%s'", cloud)
		return "", err
	}
	lastBuiltAmiId, err := utils.GetLastBuiltAmiId(packerManifestPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while getting last built AMI ID for cloud '%s'", cloud)
		return "", err
	}
	amiIdFromConfig := utils.GetAmiIdFromConfig()
	pickedAmiId := lo.Switch[string, string](amiIdFromConfig).
		Case("", lastBuiltAmiId).
		Default(amiIdFromConfig)

	publicIp, err := utils.GetPublicIp()
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while getting public IP for cloud '%s'", cloud)
		log.Print(err)
	}
	pickedIp := lo.Switch[string, string](publicIp).
		CaseF("", func() string {
			log.Printf("Using default value %s for public IP", utils.GetDefaultIp())
			return utils.MaskIp(utils.GetDefaultIp(), 32)
		}).Default(publicIp)

	generalTerraformEnv := templates_terraform_general.VarsFile{
		AMI_ID:     pickedAmiId,
		ALLOWED_IP: pickedIp,
	}

	_, err = templates_terraform_general.CraftGeneralTerraformTfVarsFile(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting general Terraform Vars file for cloud '%s'", cloud)
		return "", err
	}

	switch cloud {
	case "aws":
		cloudTfVarsFilePath, err = templates_terraform_aws.CraftAWS_TerraformTfVarsFile(cloud)
		if err != nil {
			err = errors.Wrapf(err, "Error occurred while crafting AWS Terraform Vars file for cloud '%s'", cloud)
			return "", err
		}
		return cloudTfVarsFilePath, nil
	default:
		err = errors.Errorf("Cloud '%s' is not supported", cloud)
		return "", err
	}
}
