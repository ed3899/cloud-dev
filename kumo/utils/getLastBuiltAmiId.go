package utils

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type PackerBuild struct {
	PackerRunUUID string `json:"packer_run_uuid"`
	ArtifactId    string `json:"artifact_id"`
}

type PackerManifest struct {
	Builds []*PackerBuild
	LastRunUUID string `json:"last_run_uuid"`
}

func GetLastBuiltAmiId(packerManifestPath string) (amiId string, err error) {
	file, err := os.Open(packerManifestPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while opening packer manifest file '%s'", packerManifestPath)
		return "", err
	}
	defer file.Close()

	// Decode packer manifest
	decoder := json.NewDecoder(file)
	var packerManifest PackerManifest
	err = decoder.Decode(&packerManifest)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while decoding packer manifest file '%s'", packerManifestPath)
		return "", err
	}

	// Get AMI ID for last Packer build
	lastBuildAMI_Id := lo.Filter(packerManifest.Builds, func(pb *PackerBuild, index int) bool {
		return pb.PackerRunUUID == packerManifest.LastRunUUID
	})

	if len(lastBuildAMI_Id) == 0 {
		err = errors.New("No AMI ID found for last Packer build")
		return "", err
	}

	return lastBuildAMI_Id[0].ArtifactId, nil
}