package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

// Retrieves the packer manifest file absolute path for the given cloud. If the manifest is not
// present, it returns an error.
//
// This function should be ran or called from the root directory of the project. Otherwise, it will
// fail to find the manifest file
//
// Example:
//
//	("aws") -> ("/home/kumo/packer/aws/manifest.json", nil)
func GetPackerManifestPathTo(cloud string) (packerManifestAbsPath string, err error) {
	packerManifestAbsPath, err = filepath.Abs(filepath.Join("packer", cloud, "manifest.json"))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to packer manifest for cloud '%s'", cloud)
		return "", err
	}

	if FileNotPresent(packerManifestAbsPath) {
		err = errors.Errorf("Packer manifest file not found at '%s'", packerManifestAbsPath)
		return "", err
	}

	return
}
