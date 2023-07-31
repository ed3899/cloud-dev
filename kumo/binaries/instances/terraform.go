package instances

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Terraform struct {
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *download.Zip
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

	if absPathToExecutable, err = filepath.Abs(filepath.Join(dirs.DEPENDENCIES_DIR_NAME, NAME, executableName)); err != nil {
		err = oopsBuilder.
			With("dirs.DEPENDENCIES_DIR_NAME", dirs.DEPENDENCIES_DIR_NAME).
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

	if zipPath, err = filepath.Abs(filepath.Join(dirs.DEPENDENCIES_DIR_NAME, NAME, zipName)); err != nil {
		err = oopsBuilder.
			With("dirs.DEPENDENCIES_DIR_NAME", dirs.DEPENDENCIES_DIR_NAME).
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
		AbsPathToExecutable: absPathToExecutable,
		AbsPathToRunDir:     absPathToRunDir,
		Zip: &download.Zip{
			Name:          zipName,
			AbsPath:       zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}

	return
}

func (t *Terraform) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(t.AbsPathToExecutable)
}

func (t *Terraform) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(t.AbsPathToExecutable)
}

func (t *Terraform) Init() (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "init")
		oopsBuilder = oops.
				Code("terraform_init_failed")

		cmdErr error
	)

	if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
		err = oopsBuilder.
			Wrapf(cmdErr, "Error occured while initializing terraform")
		return
	}

	return
}

func (t *Terraform) Up() (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "apply", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_up_failed")

		cmdErr error
	)

	if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
		err = oopsBuilder.
			Wrapf(cmdErr, "Error occured while deploying terraform resources")
		return
	}

	return
}

func (t *Terraform) Destroy() (err error) {
	var (
		cmd         = exec.Command(t.AbsPathToExecutable, "destroy", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_destroy_failed")

		cmdErr error
	)

	if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
		err = oopsBuilder.
			Wrapf(cmdErr, "Error occured while destroying terraform resources")
		return
	}

	return
}
