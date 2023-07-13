package templates

import (
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func GetAmiToBeUsed(packerManifestPath, cloud string) (amiIdToBeUsed string, err error) {
	// Get last built AMI ID
	lastBuiltAmiId, err := utils.GetLastBuiltAmiId(packerManifestPath)
	var lbaierr error
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while getting last built AMI ID for cloud '%s'", cloud)
		lbaierr = err
	}

	// Set AMI ID to be used based on config file
	kumoConfigAmiId := viper.GetString("Up.AMI_Id")
	switch {
	// If AMI ID is not set in config file and no AMI ID is found in manifest.json, return error
	case lastBuiltAmiId == "" && kumoConfigAmiId == "":
		totalError := errors.Wrap(lbaierr, "No AMI ID is set in config file and no AMI ID is found in manifest.json. Please make sure that the manifest.json generated after a packer build is still there and intact. By default, kumo checks for the last built artifact id and uses that ami id for deploying. If you have removed that from AMI the aws console then make sure to either generate a new build or fill in the config with a specific AMI id. Please don't tamper with the manifest.json, otherwise you will have to either input an id or rebuild a new one.")
		return "", totalError
	// If AMI ID is not set in config file, use the last built AMI ID
	case kumoConfigAmiId == "":
		amiIdToBeUsed = lastBuiltAmiId
		// If AMI ID is set in config file, use the AMI ID from config file
	case kumoConfigAmiId != "":
		amiIdToBeUsed = kumoConfigAmiId
	}

	return amiIdToBeUsed, nil
}