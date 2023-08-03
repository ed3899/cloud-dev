package terraform

import (
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Instance struct {
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *download.Zip
}

func NewInstance(config binaries.ConfigI) (instance *Instance, err error) {
	var (
		dependenciesDirName     = config.GetDependenciesDirName()
		terraformName           = config.GetToolName()
		terraformDirName        = terraformName
		terraformVersion        = config.GetToolVersion()
		terraformExecutableName = config.GetToolExecutableName()
		terraformZipName        = config.GetToolZipName()
		currentOs, currentArch  = config.GetCurrentHostSpecs()
		terraformUrl            = config.CreateHashicorpURL(terraformName, terraformVersion, currentOs, currentArch)
		oopsBuilder             = oops.
					Code("new_terraform_failed")

		absPathToKumoExecutable      string
		absPathToKumoExecutableDir   string
		absPathToTerraformExecutable string
		absPathToTerraformRunDir     string
		absPathToTerraformZip        string
		terraformZipContentLength    int64
	)

	if absPathToKumoExecutable, err = config.GetKumoExecutableAbsPath(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to kumo executable")
		return
	}

	absPathToKumoExecutableDir = filepath.Dir(absPathToKumoExecutable)
	absPathToTerraformExecutable = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, terraformDirName, terraformExecutableName)
	absPathToTerraformRunDir = filepath.Join(absPathToKumoExecutableDir, terraformDirName)
	absPathToTerraformZip = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, terraformDirName, terraformZipName)

	if terraformZipContentLength, err = config.GetContentLength(terraformUrl); err != nil {
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
		oopsBuilder = oops.
				Code("terraform_init_failed")
		cmd         = exec.Command(i.AbsPathToExecutable, "init")
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
