package templates

import (
	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	templates_terraform_general "github.com/ed3899/kumo/templates/terraform/general"
	"github.com/pkg/errors"
)

func CraftCloudTerraformTfVarsFile(cloud string) (cloudTfVarsFilePath string, err error) {
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
