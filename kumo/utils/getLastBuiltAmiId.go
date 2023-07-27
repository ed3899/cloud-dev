package utils

import (
	"encoding/json"
	"log"
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
func GetLastBuiltAmiId(packerManifestAbsPath string) (amiId string, err error) {
	var (
		packerManifest = new(PackerManifest)
		oopsBuilder    = oops.Code("get_last_built_ami_id_failed").
				With("packerManifestAbsPath", packerManifestAbsPath)

		lastBuildArtifact  []*PackerBuild
		packerManifestFile *os.File
		jsonDecoder        *json.Decoder
	)

	// Check if packer manifest path is absolute
	if !filepath.IsAbs(packerManifestAbsPath) {
		err = oopsBuilder.
			Wrapf(err, "The packer manifest path is not absolute: '%s'", packerManifestAbsPath)
		return
	}

	// Open packer manifest file
	if packerManifestFile, err = os.Open(packerManifestAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while opening packer manifest file '%s'", packerManifestAbsPath)
		return
	}
	defer func() {
		if err := packerManifestFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "Error occurred while closing packer manifest file: '%s'", packerManifestFile.Name()),
			)
		}
	}()

	// Decode packer manifest
	jsonDecoder = json.NewDecoder(packerManifestFile)
	if err = jsonDecoder.Decode(packerManifest); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while decoding packer manifest file '%s'", packerManifestAbsPath)
		return
	}

	// Get last built artifact id for last Packer build
	lastBuildArtifact = lo.Filter(packerManifest.Builds, func(pb *PackerBuild, index int) bool {
		return pb.PackerRunUUID == packerManifest.LastRunUUID
	})

	if len(lastBuildArtifact) == 0 {
		err = oopsBuilder.
			With("lastBuildArtifact", lastBuildArtifact).
			Wrapf(err, "No AMI ID found for last Packer build")
		return
	}

	// Extract only the AMI ID
	amiId = strings.Split(lastBuildArtifact[0].ArtifactId, ":")[1]

	return
}
