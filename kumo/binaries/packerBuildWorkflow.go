package binaries

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func PackerBuildWorkflow() (err error) {
	var (
		cloud                    Cloud
		packer                   *Packer
		absPathToInitialLocation string
		absPathToCloudRunDir     string
		varsFile                 *HashicorpVars
	)

	// A. Instantiate packer
	if packer, err = NewPacker(); err != nil {
		err = errors.Wrap(err, "Error occurred while creating new packer")
		return
	}

	// A. Download and extract if not installed
	if packer.IsNotInstalled() {
		err = DownloadAndExtractWorkflow(packer.Zip, filepath.Dir(packer.AbsPathToExecutable))
		if err != nil {
			err = errors.Wrap(err, "Error occurred while downloading and extracting packer")
			return
		}
	}

	// B. Set cloud
	if cloud, err = GetCloud(); err != nil {
		err = errors.Wrap(err, "Error occurred while getting cloud")
		return
	}

	// C. Instantiate vars file
	if varsFile, err = NewHashicorpVars(packer.ID, cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while instantiating hashicorp vars")
		return
	}

	// C. Create vars file
	if err = varsFile.Create(); err != nil {
		err = errors.Wrap(err, "Error occurred while creating vars file")
		return
	}

	// D. Get abs path to cloud run directory
	if absPathToCloudRunDir, err = packer.GetAbsPathToCloudRunDir(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while getting absolute path to cloud run directory")
		return
	}

	// D. Set initial location
	if absPathToInitialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return
	}

	// D. Change directory to packer run directory and defer changing back to initial location
	if err = os.Chdir(absPathToCloudRunDir); err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to packer run directory")
		return
	}
	defer func() {
		if err = os.Chdir(absPathToInitialLocation); err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory back to initial location")
			return
		}
	}()

	// E. Set plugin path
	if err = packer.SetPluginPath(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while setting plugin path")
		return
	}
	defer func() {
		if err = packer.UnsetPluginPath(cloud); err != nil {
			err = errors.Wrap(err, "Error occurred while unsetting plugin path")
			return
		}
	}()

	// F. Initialize
	if err = packer.Init(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while initializing packer")
		return
	}

	// G. Build
	if err = packer.Build(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while building packer")
		return
	}

	return
}
