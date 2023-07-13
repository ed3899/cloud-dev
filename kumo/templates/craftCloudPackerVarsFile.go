package templates

import "github.com/pkg/errors"

func CraftCloudPackerVarsFile(cloud string) (cloudPackerVarsFile string, err error) {
	_, err = CraftGeneralPackerVarsFile(cloud)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing general Packer vars file")
		return "", err
	}
	
	switch cloud {
	case "aws":
		cloudPackerVarsFile, err = CraftAWSPackerVarsFile()
		if err != nil {
			err = errors.Wrap(err, "Error occurred while crafting AWS Packer Vars file")
			return "", err
		}
		return cloudPackerVarsFile, nil
	default:
		err = errors.Errorf("Cloud template for '%s' not supported", cloud)
		return "", err
	}
}
