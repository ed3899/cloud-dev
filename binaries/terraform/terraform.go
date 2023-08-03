package terraform

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Instance struct {
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *download.Zip
}

func NewInstance() (instance *Instance, err error) {
	var (
		dependenciesDirName     = dirs.DEPENDENCIES_DIR_NAME
		terraformName           = tool.TERRAFORM_NAME
		terraformDirName        = terraformName
		terraformVersion        = tool.TERRAFORM_VERSION
		terraformExecutableName = fmt.Sprintf("%s.exe", terraformName)
		terraformZipName        = fmt.Sprintf("%s.zip", terraformName)
		os, arch                = utils.GetCurrentHostSpecs()
		terraformUrl            = utils.CreateHashicorpURL(terraformName, terraformVersion, os, arch)
		oopsBuilder             = oops.
					Code("new_terraform_failed")

		absPathToTerraformExecutable string
		absPathToTerraformRunDir     string
		absPathToTerraformZip        string
		terraformZipContentLength    int64
	)

	if absPathToTerraformExecutable, err = filepath.Abs(filepath.Join(dependenciesDirName, terraformDirName, terraformExecutableName)); err != nil {
		err = oopsBuilder.
			With("dependenciesDirName", dependenciesDirName).
			With("terraformDirName", terraformDirName).
			Wrapf(err, "failed to create absolute path to: %s", terraformExecutableName)
		return
	}

	if absPathToTerraformRunDir, err = filepath.Abs(terraformDirName); err != nil {
		err = oopsBuilder.
			With("terraformDirName", terraformDirName).
			Wrapf(err, "failed to create absolute path to run dir")
		return
	}

	if absPathToTerraformZip, err = filepath.Abs(filepath.Join(dependenciesDirName, terraformName, terraformZipName)); err != nil {
		err = oopsBuilder.
			With("dependenciesDirName", dependenciesDirName).
			With("terraformName", terraformName).
			Wrapf(err, "failed to create absolute path to: %s", terraformZipName)
		return
	}

	if terraformZipContentLength, err = utils.GetContentLength(terraformUrl); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", terraformUrl)
		return
	}

	instance = &Instance{
		AbsPathToExecutable: absPathToTerraformExecutable,
		AbsPathToRunDir:     absPathToTerraformRunDir,
		Zip: &download.Zip{
			Name:          terraformZipName,
			AbsPath:       absPathToTerraformZip,
			URL:           terraformUrl,
			ContentLength: terraformZipContentLength,
		},
	}

	return
}

func (i *Instance) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(i.AbsPathToExecutable)
}

func (i *Instance) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(i.AbsPathToExecutable)
}

func (i *Instance) Init() (err error) {
	var (
		cmd         = exec.Command(i.AbsPathToExecutable, "init")
		oopsBuilder = oops.
				Code("terraform_init_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform init command")
		return
	}

	return
}

func (i *Instance) Up() (err error) {
	var (
		cmd         = exec.Command(i.AbsPathToExecutable, "apply", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_up_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform apply command")
		return
	}

	return
}

func (i *Instance) Destroy() (err error) {
	var (
		cmd         = exec.Command(i.AbsPathToExecutable, "destroy", "-auto-approve")
		oopsBuilder = oops.
				Code("terraform_destroy_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming terraform destroy command")
		return
	}

	return
}
