package packer_manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"
)

type PackerBuild struct {
	PackerRunUUID string `json:"packer_run_uuid"`
	ArtifactId    string `json:"artifact_id"`
}

type PackerManifest struct {
	Builds      []*PackerBuild
	LastRunUUID string `json:"last_run_uuid"`
}

// Returns the AMI ID of the last built AMI. The AMI ID is extracted from the Packer manifest file.
//
// Example:
//
//	("packer/aws/manifest.json") -> ("ami-0c3fd0f5d33134a76", nil)
func GetLastBuiltAmiIdFromPackerManifest(
	packerManifestAbsPath string,
) (string, error) {
	packerManifest := &PackerManifest{}
	oopsBuilder := oops.
		Code("get_last_built_ami_id_failed").
		With("packerManifestAbsPath", packerManifestAbsPath)

	// Check if packer manifest path is absolute
	if !filepath.IsAbs(packerManifestAbsPath) {
		err := oopsBuilder.
			Errorf("The packer manifest path is not absolute: '%s'", packerManifestAbsPath)

		return "", err
	}

	// Open packer manifest file
	packerManifestFile, err := os.Open(packerManifestAbsPath)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while opening packer manifest file '%s'", packerManifestAbsPath)

		return "", err
	}
	defer packerManifestFile.Close()

	// Decode packer manifest
	jsonDecoder := json.NewDecoder(packerManifestFile)
	err = jsonDecoder.Decode(packerManifest)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while decoding packer manifest file '%s'", packerManifestAbsPath)

		return "", err
	}

	// Get last built artifact id for last Packer build
	lastBuildArtifact := lo.Filter(packerManifest.Builds, func(pb *PackerBuild, index int) bool {
		return pb.PackerRunUUID == packerManifest.LastRunUUID
	})

	if len(lastBuildArtifact) == 0 {
		err := oopsBuilder.
			With("lastBuildArtifact", lastBuildArtifact).
			Wrapf(err, "No AMI ID found for last Packer build")
		return "", err
	}

	// Extract only the AMI ID
	return strings.Split(lastBuildArtifact[0].ArtifactId, ":")[1], nil
}
