package templates

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type GeneralTerraformEnvironment struct {
	*AWS_TerraformEnvironment
	AMI_ID     string
	ALLOWED_IP string
}



func CraftGeneralTerraformTfVarsFile(cloud string) (generalTerraformVarsPath string, err error) {
	// Get packer manifest
	packerManifestPath, err := utils.CraftAbsolutePath("packer", cloud, "manifest.json")
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to packer manifest for cloud '%s'", cloud)
		return "", err
	}

	// Get AMI ID to be used
	amiIdToBeUsed, err := GetAmiToBeUsed(packerManifestPath, cloud)
	if err != nil {
		log.Fatal(err)
	}
	// Get host public IP
	hostPublicIp := viper.GetString("ALLOWED_ID")

	genEnv := &GeneralTerraformEnvironment{
		AMI_ID: amiIdToBeUsed,
		ALLOWED_IP: hostPublicIp,
	}

	generalTerraformVarsPath, err = CraftGenericTerraformTfVarsFile[*GeneralTerraformEnvironment](cloud, "GeneralTerraformTfVars.tmpl", "general.auto.tfvars", genEnv)

	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting generic Terraform Vars file")
		return "", err
	}

	return generalTerraformVarsPath, nil

}
