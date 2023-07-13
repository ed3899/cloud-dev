package templates

import "github.com/pkg/errors"

import (
	templates_packer_general "github.com/ed3899/kumo/templates/packer/general"
	templates_packer_aws "github.com/ed3899/kumo/templates/packer/aws"
)

func CraftCloudPackerVarsFile(cloud string) (cloudPackerVarsFilePath string, err error) {
	_, err = templates_packer_general.CraftGeneralPackerVarsFile(cloud)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing general Packer vars file")
		return "", err
	}

	switch cloud {
	case "aws":
		cloudPackerVarsFilePath, err = templates_packer_aws.CraftAWSPackerVarsFile()
		if err != nil {
			err = errors.Wrap(err, "Error occurred while crafting AWS Packer Vars file")
			return "", err
		}
		return cloudPackerVarsFilePath, nil
	default:
		err = errors.Errorf("Cloud template for '%s' not supported", cloud)
		return "", err
	}
}
