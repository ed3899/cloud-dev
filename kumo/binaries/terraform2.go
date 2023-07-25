package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type TerraformI interface {
	Init(Cloud) error
	Up(Cloud) error
	Destroy(Cloud) error
}

type Terraform struct {
	ID                  Tool
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *Zip
}

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
)

func NewTerraform() (terraform *Terraform, err error) {
	const (
		name    = "terraform"
		version = "0.14.7"
	)

	var (
		executableName      = fmt.Sprintf("%s.exe", name)
		zipName             = fmt.Sprintf("%s.zip", name)
		os, arch            = utils.GetCurrentHostSpecs()
		url                 = utils.CreateHashicorpURL(name, version, os, arch)
		depDirName          = utils.GetDependenciesDirName()
		absPathToExecutable string
		absPathToRunDir     string
		contentLength       int64
		zipPath             string
	)

	if absPathToExecutable, err = filepath.Abs(filepath.Join(depDirName, name, executableName)); err != nil {
		err = errors.Wrapf(err, "failed to create executable path to: %s", executableName)
		return
	}

	if absPathToRunDir, err = filepath.Abs(name); err != nil {
		err = errors.Wrapf(err, "failed to create run path to: %s", name)
		return
	}

	if zipPath, err = filepath.Abs(filepath.Join(depDirName, name, zipName)); err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %s", zipName)
		return
	}

	if contentLength, err = utils.GetContentLength(url); err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	terraform = &Terraform{
		ID:                  TerraformID,
		AbsPathToExecutable: absPathToExecutable,
		AbsPathToRunDir:     absPathToRunDir,
		Zip: &Zip{
			Name:          zipName,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}
	return
}

func (t *Terraform) SetCloudCredentials(cloud Cloud) (err error) {
	switch cloud {
	case AWS:
		if err = os.Setenv(AWS_ACCESS_KEY_ID, viper.GetString("AWS.AccessKeyId")); err != nil {
			err = errors.Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID)
			return err
		}

		if err = os.Setenv(AWS_SECRET_ACCESS_KEY, viper.GetString("AWS.SecretAccessKey")); err != nil {
			err = errors.Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY)
			return err
		}

		return
	default:
		err = errors.Errorf("Cloud '%s' is not supported", cloud)
		return
	}
}

func (t *Terraform) UnsetCloudCredentials(cloud Cloud) (err error) {
	switch cloud {
	case AWS:
		if err = os.Unsetenv(AWS_ACCESS_KEY_ID); err != nil {
			err = errors.Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_ACCESS_KEY_ID)
			return err
		}

		if err = os.Unsetenv(AWS_SECRET_ACCESS_KEY); err != nil {
			err = errors.Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_SECRET_ACCESS_KEY)
			return err
		}

		return
	default:
		err = errors.Errorf("Cloud '%s' is not supported", cloud)
		return
	}
}

func (t *Terraform) Init(cloud Cloud) (err error) {
	var (
		cmd             = exec.Command(t.AbsPathToExecutable, "init", ".")
		cmdErr          error
		initialLocation string
	)

	// Store current working directory
	if initialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch cloud {
	case AWS:
		// Change directory to where packer will be run
		absPathToRunLocation := filepath.Join(t.AbsPathToRunDir, AWS_SUBDIR_NAME)
		if err = os.Chdir(absPathToRunLocation); err != nil {
			err = errors.Wrapf(err, "Error occurred while changing directory to %s", absPathToRunLocation)
			return
		}
		defer os.Chdir(initialLocation)

		// Initialize
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while initializing terraform for %v", cloud)
			return
		}

		return

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
		return
	}
}

func (t *Terraform) Up(cloud Cloud) (err error) {
	var (
		cmd             = exec.Command(t.AbsPathToExecutable, "apply", "-auto-approve", ".")
		cmdErr          error
		initialLocation string
	)

	// Store current working directory
	if initialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch cloud {
	case AWS:
		// Change to run directory
		absPathToRunLocation := filepath.Join(t.AbsPathToRunDir, AWS_SUBDIR_NAME)
		if err = os.Chdir(absPathToRunLocation); err != nil {
			err = errors.Wrapf(err, "Error occurred while changing directory to %s", absPathToRunLocation)
			return err
		}
		defer os.Chdir(initialLocation)

		// Run cmd
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while deploying for %v", cloud)
			return
		}
		return

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
		return
	}
}

func (t *Terraform) Destroy(cloud Cloud) (err error) {
	var (
		cmd             = exec.Command(t.AbsPathToExecutable, "destroy", "-auto-approve", ".")
		cmdErr          error
		initialLocation string
	)

	// Store current working directory
	if initialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch cloud {
	case AWS:
		// Change to run directory
		absPathToRunLocation := filepath.Join(t.AbsPathToRunDir, AWS_SUBDIR_NAME)
		if err = os.Chdir(absPathToRunLocation); err != nil {
			err = errors.Wrapf(err, "Error occurred while changing directory to %s", absPathToRunLocation)
			return err
		}
		defer os.Chdir(initialLocation)

		// Run cmd
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while destroying for %v", cloud)
			return
		}
		return

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
		return
	}
}
