package packer

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/binaries/download"
	"github.com/ed3899/kumo/binaries/instances"
	"github.com/ed3899/kumo/binaries/workflows"
	"github.com/pkg/errors"
)

// TODO remove errors for ooops
func BuildWorkflow() (err error) {
	var (
		cloud                    binaries.Cloud
		packer                   *instances.Packer
		absPathToInitialLocation string
		absPathToCloudRunDir     string
		varsFile                 *workflows.HashicorpVars
	)

	// A. Instantiate packer
	if packer, err = instances.NewPacker(); err != nil {
		return errors.Wrap(err, "Error occurred while creating new packer")
	}

	// Download and extract if not installed
	if packer.IsNotInstalled() {
		if err = download.DownloadAndExtract(packer.Zip, filepath.Dir(packer.AbsPathToExecutable)); err != nil {
			return errors.Wrap(err, "Error occurred while downloading and extracting packer")
		}
	}

	// B. Set cloud
	if cloud, err = workflows.GetCloud(); err != nil {
		return errors.Wrap(err, "Error occurred while getting cloud")
	}

	// C. Instantiate vars file
	if varsFile, err = workflows.NewHashicorpVars(packer.ToolId, cloud); err != nil {
		return errors.Wrap(err, "Error occurred while instantiating hashicorp vars")
	}

	// Create vars file
	if err = varsFile.Create(); err != nil {
		return errors.Wrap(err, "Error occurred while creating vars file")
	}

	// D. Get abs path to cloud run directory
	if absPathToCloudRunDir, err = packer.GetAbsPathToCloudRunDir(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while getting absolute path to cloud run directory")
	}

	// Set initial location
	if absPathToInitialLocation, err = os.Getwd(); err != nil {
		return errors.Wrap(err, "Error occurred while getting current working directory")
	}

	// Change directory to packer run directory and defer changing back to initial location
	if err = os.Chdir(absPathToCloudRunDir); err != nil {
		return errors.Wrap(err, "Error occurred while changing directory to packer run directory")
	}
	defer func() {
		if chdirErr := os.Chdir(absPathToInitialLocation); chdirErr != nil {
			err = errors.Wrap(chdirErr, "Error occurred while changing directory back to initial location")
		}
	}()

	// E. Set plugin path
	if err = packer.SetPluginPath(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while setting plugin path")
	}
	defer func() {
		if unsetError := packer.UnsetPluginPath(cloud); unsetError != nil {
			err = errors.Wrap(unsetError, "Error occurred while unsetting plugin path")
		}
	}()

	// F. Initialize
	if err = packer.Init(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while initializing packer")
	}

	// G. Build
	if err = packer.Build(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while building packer")
	}

	return
}
