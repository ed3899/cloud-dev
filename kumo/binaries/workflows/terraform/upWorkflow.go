package terraform

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/binaries/download"
	"github.com/ed3899/kumo/binaries/instances"
	"github.com/ed3899/kumo/binaries/workflows"
	"github.com/pkg/errors"
)

func UpWorkflow() (err error) {
	var (
		cloud                    binaries.Cloud
		terraform                *instances.Terraform
		absPathToInitialLocation string
		absPathToCloudRunDir     string
		varsFile                 *workflows.HashicorpVars
	)

	// A. Instantiate terraform
	if terraform, err = instances.NewTerraform(); err != nil {
		return errors.Wrap(err, "Error occurred while creating new terraform")
	}

	// Download and extract if not installed
	if terraform.IsNotInstalled() {
		if err = download.DownloadAndExtract(terraform.Zip, filepath.Dir(terraform.AbsPathToExecutable)); err != nil {
			return errors.Wrap(err, "Error occurred while downloading and extracting terraform")
		}
	}

	// B. Set cloud
	if cloud, err = workflows.GetCloud(); err != nil {
		return errors.Wrap(err, "Error occurred while getting cloud")
	}

	// C. Instantiate vars file
	if varsFile, err = workflows.NewHashicorpVars(terraform.ID, cloud); err != nil {
		return errors.Wrap(err, "Error occurred while instantiating hashicorp vars")
	}

	// Create vars file
	if err = varsFile.Create(); err != nil {
		return errors.Wrap(err, "Error occurred while creating vars file")
	}

	// D. Get abs path to cloud run directory
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
		if chDirError := os.Chdir(absPathToInitialLocation); chDirError != nil {
			err = errors.Wrap(chDirError, "Error occurred while changing directory back to initial location")
		}
	}()

	// E. Set cloud credentials and defer unsetting
	if err = terraform.SetCloudCredentials(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while setting cloud credentials")
	}
	defer func() {
		if unsetCloudErr := terraform.UnsetCloudCredentials(cloud); unsetCloudErr != nil {
			err = errors.Wrap(unsetCloudErr, "Error occurred while unsetting cloud credentials")
		}
	}()

	// F. Initialize
	if err = terraform.Init(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while initializing terraform")
	}

	// G. Up
	if err = terraform.Up(cloud); err != nil {
		return errors.Wrap(err, "Error occurred while running terraform up")
	}

	return
}
