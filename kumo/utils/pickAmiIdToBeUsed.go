package utils

import "github.com/pkg/errors"

// Return the AMI ID to be used for deployment. If the lastBuildAmiId is empty, it will return an error.
// If the amiIdFromConfig is empty, it will return the lastBuildAmiId. Otherwise, it will return the
// amiIdFromConfig.
func PickAmiIdToBeUsed(lastBuildAmiId, amiIdFromConfig string) (amiIdToBeUsed string, err error) {
	switch {
	case lastBuildAmiId == "":
		err = errors.New("The lastBuildAmiId is not present. Please make sure that the manifest.json generated after a packer build is still there and intact. By default, kumo checks for the last built artifact id and uses that ami id for deploying. Please don't tamper with the manifest.json, otherwise you will have to remove your AMI from the AWS Console, perform a cleanup and build up a new one.")
	case amiIdFromConfig == "":
		amiIdToBeUsed = lastBuildAmiId
	case amiIdFromConfig != "":
		amiIdToBeUsed = amiIdFromConfig
	default:
		err = errors.New("Something went wrong while picking the AMI ID to be used for deployment")
	}

	return
}
