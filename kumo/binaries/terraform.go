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
		NAME    = "terraform"
		VERSION = "1.5.3"
	)

	var (
		executableName      = fmt.Sprintf("%s.exe", NAME)
		zipName             = fmt.Sprintf("%s.zip", NAME)
		os, arch            = utils.GetCurrentHostSpecs()
		url                 = utils.CreateHashicorpURL(NAME, VERSION, os, arch)
		absPathToExecutable string
		absPathToRunDir     string
		contentLength       int64
		zipPath             string
	)

	if absPathToExecutable, err = filepath.Abs(filepath.Join(DEPENDENCIES_DIR_NAME, NAME, executableName)); err != nil {
		err = errors.Wrapf(err, "failed to create executable path to: %s", executableName)
		return
	}

	if absPathToRunDir, err = filepath.Abs(NAME); err != nil {
		err = errors.Wrapf(err, "failed to create run path to: %s", NAME)
		return
	}

	if zipPath, err = filepath.Abs(filepath.Join(DEPENDENCIES_DIR_NAME, NAME, zipName)); err != nil {
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
			AbsPath:       zipPath,
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

	default:
		err = errors.Errorf("Cloud '%v' is not supported", cloud)
	}

	return
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

	default:
		err = errors.Errorf("Cloud '%v' is not supported", cloud)
	}

	return
}

func (t *Terraform) GetAbsPathToCloudRunDir(cloud Cloud) (cloudRunDir string, err error) {
	switch cloud {
	case AWS:
		cloudRunDir = filepath.Join(t.AbsPathToRunDir, AWS_SUBDIR_NAME)
	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}
	return
}

func (t *Terraform) IsInstalled() bool {
	return utils.FilePresent(t.AbsPathToExecutable)
}

func (t *Terraform) IsNotInstalled() bool {
	return utils.FileNotPresent(t.AbsPathToExecutable)
}

func (t *Terraform) Init(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(t.AbsPathToExecutable, "init", ".")
		cmdErr error
	)

	switch cloud {
	case AWS:
		// Run cmd
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while initializing terraform for %v", cloud)
			return
		}

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}

	return
}

func (t *Terraform) Up(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(t.AbsPathToExecutable, "apply", "-auto-approve", ".")
		cmdErr error
	)

	switch cloud {
	case AWS:
		// Run cmd
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while deploying for %v", cloud)
			return
		}

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}

	return
}

func (t *Terraform) Destroy(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(t.AbsPathToExecutable, "destroy", "-auto-approve", ".")
		cmdErr error
	)

	switch cloud {
	case AWS:
		// Run cmd
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while destroying for %v", cloud)
			return
		}

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}

	return
}
