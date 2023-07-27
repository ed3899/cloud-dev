package instances

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type TerraformI interface {
	Init(binaries.Cloud) error
	Up(binaries.Cloud) error
	Destroy(binaries.Cloud) error
}

type Terraform struct {
	ID                  binaries.Tool
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
		executableName = fmt.Sprintf("%s.exe", NAME)
		zipName        = fmt.Sprintf("%s.zip", NAME)
		os, arch       = utils.GetCurrentHostSpecs()
		url            = utils.CreateHashicorpURL(NAME, VERSION, os, arch)
		oopsBuilder    = oops.
				Code("new_terraform_failed")

		absPathToExecutable string
		absPathToRunDir     string
		contentLength       int64
		zipPath             string
	)

	if absPathToExecutable, err = filepath.Abs(filepath.Join(binaries.DEPENDENCIES_DIR_NAME, NAME, executableName)); err != nil {
		err = oopsBuilder.
			With("DEPENDENCIES_DIR_NAME", binaries.DEPENDENCIES_DIR_NAME).
			With("NAME", NAME).
			Wrapf(err, "failed to create absolute path to: %s", executableName)
		return
	}

	if absPathToRunDir, err = filepath.Abs(NAME); err != nil {
		err = oopsBuilder.
			With("NAME", NAME).
			Wrapf(err, "failed to create absolute path to run dir")
		return
	}

	if zipPath, err = filepath.Abs(filepath.Join(binaries.DEPENDENCIES_DIR_NAME, NAME, zipName)); err != nil {
		err = oopsBuilder.
			With("DEPENDENCIES_DIR_NAME", binaries.DEPENDENCIES_DIR_NAME).
			With("NAME", NAME).
			Wrapf(err, "failed to create absolute path to: %s", zipName)
		return
	}

	if contentLength, err = utils.GetContentLength(url); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	terraform = &Terraform{
		ID:                  binaries.TerraformID,
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

func (t *Terraform) SetCloudCredentials(cloud binaries.Cloud) (err error) {
	var (
		oopsBuilder = oops.
			Code("terraform_set_cloud_credentials_failed").
			With("cloud", cloud)
	)

	switch cloud {
	case binaries.AWS:
		if err = os.Setenv(AWS_ACCESS_KEY_ID, viper.GetString("AWS.AccessKeyId")); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while setting %s environment variable", AWS_ACCESS_KEY_ID)
			return
		}

		if err = os.Setenv(AWS_SECRET_ACCESS_KEY, viper.GetString("AWS.SecretAccessKey")); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while setting %s environment variable", AWS_SECRET_ACCESS_KEY)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' is not supported", cloud)
		return
	}

	return
}

func (t *Terraform) UnsetCloudCredentials(cloud binaries.Cloud) (err error) {
	var (
		oopsBuilder = oops.
			Code("terraform_unset_cloud_credentials_failed").
			With("cloud", cloud)
	)

	switch cloud {
	case binaries.AWS:
		if err = os.Unsetenv(AWS_ACCESS_KEY_ID); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_ACCESS_KEY_ID)
			return
		}

		if err = os.Unsetenv(AWS_SECRET_ACCESS_KEY); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while unsetting %s environment variable", AWS_SECRET_ACCESS_KEY)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' is not supported", cloud)
		return
	}

	return
}

func (t *Terraform) GetAbsPathToCloudRunDir(cloud binaries.Cloud) (cloudRunDir string, err error) {
	var (
		oopsBuilder = oops.
			Code("terraform_get_abs_path_to_cloud_run_dir_failed").
			With("cloud", cloud)
	)

	switch cloud {
	case binaries.AWS:
		cloudRunDir = filepath.Join(t.AbsPathToRunDir, binaries.AWS_SUBDIR_NAME)

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' is not supported", cloud)
		return
	}

	return
}

func (t *Terraform) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(t.AbsPathToExecutable)
}

func (t *Terraform) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(t.AbsPathToExecutable)
}

func (t *Terraform) Init(cloud binaries.Cloud) (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "init")
		oopsBuilder = oops.
				Code("terraform_init_failed").
				With("cloud", cloud)

		cmdErr error
	)

	switch cloud {
	case binaries.AWS:
		if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
			err = oopsBuilder.
				Wrapf(cmdErr, "Error occured while initializing terraform for %v", cloud)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}

func (t *Terraform) Up(cloud binaries.Cloud) (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "apply", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_up_failed").
				With("cloud", cloud)

		cmdErr error
	)

	switch cloud {
	case binaries.AWS:
		if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
			err = oopsBuilder.
				Wrapf(cmdErr, "Error occured while deploying for %v", cloud)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}

func (t *Terraform) Destroy(cloud binaries.Cloud) (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "destroy", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_destroy_failed").
				With("cloud", cloud)

		cmdErr error
	)

	switch cloud {
	case binaries.AWS:
		if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
			err = oopsBuilder.
				Wrapf(cmdErr, "Error occured while destroying for %v", cloud)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}
