package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type PackerBuild struct {
	PackerRunUUID string `json:"packer_run_uuid"`
	ArtifactId    string `json:"artifact_id"`
}

type PackerManifest struct {
	Builds      []*PackerBuild
	LastRunUUID string `json:"last_run_uuid"`
}

// Returns the AMI ID of the last built AMI. The AMI ID is extracted from the Packer manifest file. This function
// expects the absolute path to the Packer manifest file.
//
// Example:
//
//	("packer/aws/manifest.json") -> ("ami-0c3fd0f5d33134a76", nil)
func GetLastBuiltAmiId(packerManifestAbsPath string) (amiId string, err error) {
	var (
		lastBuildArtifact []*PackerBuild
		packerManifestFile  *os.File
		jsonDecoder         *json.Decoder
		
		packerManifest      = new(PackerManifest)
	)

	// Check if packer manifest path is absolute
	if !filepath.IsAbs(packerManifestAbsPath) {
		return "", errors.New("The packer manifest path is not absolute")
	}

	// Open packer manifest file
	if packerManifestFile, err = os.Open(packerManifestAbsPath); err != nil {
		return "", errors.Wrapf(err, "Error occurred while opening packer manifest file '%s'", packerManifestAbsPath)
	}
	defer func() {
		if errClosingPackerManifest := packerManifestFile.Close(); errClosingPackerManifest != nil {
			err = errors.Wrapf(errClosingPackerManifest, "Error occurred while closing packer manifest file '%s'", packerManifestAbsPath)
		}
	}()

	// Decode packer manifest
	jsonDecoder = json.NewDecoder(packerManifestFile)
	if err = jsonDecoder.Decode(packerManifest); err != nil {
		return "", errors.Wrapf(err, "Error occurred while decoding packer manifest file '%s'", packerManifestAbsPath)
	}

	// Get last built artifact id for last Packer build
	lastBuildArtifact = lo.Filter(packerManifest.Builds, func(pb *PackerBuild, index int) bool {
		return pb.PackerRunUUID == packerManifest.LastRunUUID
	})

	if len(lastBuildArtifact) == 0 {
		return "", errors.New("No AMI ID found for last Packer build")
	}

	// Extract only the AMI ID
	amiId = strings.Split(lastBuildArtifact[0].ArtifactId, ":")[1]

	return
}
