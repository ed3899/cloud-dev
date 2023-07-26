package binaries

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func TerraformDestroyWorkflow() (err error) {
	var (
		cloud                    Cloud
		terraform                *Terraform
		absPathToInitialLocation string
		absPathToCloudRunDir     string
	)

	// A. Instantiate terraform
	if terraform, err = NewTerraform(); err != nil {
		return errors.Wrap(err, "Error occurred while creating new terraform")
	}

	// Download and extract if not installed
	if terraform.IsNotInstalled() {
		if err = DownloadAndExtractWorkflow(terraform.Zip, filepath.Dir(terraform.AbsPathToExecutable)); err != nil {
			return errors.Wrap(err, "Error occurred while downloading and extracting terraform")
		}
	}

	// B. Set cloud
	if cloud, err = GetCloud(); err != nil {
		return errors.Wrap(err, "Error occurred while getting cloud")
	}

	// C. Get abs path to cloud run directory
	if absPathToCloudRunDir, err = terraform.GetAbsPathToCloudRunDir(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while getting absolute path to cloud run directory")
	}

	// Set initial location
	if absPathToInitialLocation, err = os.Getwd(); err != nil {
		return errors.Wrap(err, "Error occurred while getting current working directory")
	}

	// Change directory to terraform run directory and defer changing back to initial location
	if err = os.Chdir(absPathToCloudRunDir); err != nil {
		return errors.Wrap(err, "Error occurred while changing directory to packer run directory")
	}
	defer func() {
		if err = os.Chdir(absPathToInitialLocation); err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory back to initial location")
		}
	}()

	// D. Set cloud credentials and defer unsetting
	if err = terraform.SetCloudCredentials(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while setting cloud credentials")
	}
	defer func() {
		if err = terraform.UnsetCloudCredentials(cloud); err != nil {
			err = errors.Wrap(err, "Error occurred while unsetting cloud credentials")
		}
	}()

	// E. Initialize
	if err = terraform.Init(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while initializing terraform")
	}

	// F. Destroy
	if err = terraform.Destroy(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while destroying")
	}

	return
}
