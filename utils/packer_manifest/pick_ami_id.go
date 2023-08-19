package packer_manifest

import (
	"github.com/samber/oops"
)

// Return the AMI ID to be used for deployment. If the lastBuildAmiId is empty, it will return an error.
// If the amiIdFromConfig is empty, it will return the lastBuildAmiId. Otherwise, it will return the
// amiIdFromConfig.
func PickAmiId(
	lastBuildAmiId,
	amiIdFromConfig string,
) (string, error) {
	oopsBuilder := oops.
		Code("PickAmiId").
		In("utils").
		In("packer_manifest").
		With("lastBuildAmiId", lastBuildAmiId).
		With("amiIdFromConfig", amiIdFromConfig)

	switch {
	case lastBuildAmiId == "":
		err := oopsBuilder.Errorf("Last ami id not provided. Please make sure that the manifest.json generated after a packer build is still there and intact. By default, kumo checks for the last built artifact id and uses that ami id for deploying. Please don't tamper with the manifest.json, otherwise you will have to remove your AMI from the AWS Console, perform a cleanup and build up a new one.")

		return "", err

	case amiIdFromConfig == "":
		amiIdToBeUsed := lastBuildAmiId

		return amiIdToBeUsed, nil

	case amiIdFromConfig != "":
		amiIdToBeUsed := amiIdFromConfig

		return amiIdToBeUsed, nil

	default:
		err := oopsBuilder.Errorf("Something went wrong while picking the AMI ID to be used for deployment")

		return "", err
	}
}
