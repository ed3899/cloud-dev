package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/download"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Packer struct {
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *download.Zip
}

func NewPacker() (packer *Packer, err error) {
	var (
		dependenciesDirName  = dirs.DEPENDENCIES_DIR_NAME
		packerName           = tool.PACKER_NAME
		packerDirName        = packerName
		packerVersion        = tool.PACKER_VERSION
		packerExecutableName = fmt.Sprintf("%s.exe", packerName)
		packerZipName        = fmt.Sprintf("%s.zip", packerName)
		os, arch             = utils.GetCurrentHostSpecs()
		packerUrl            = utils.CreateHashicorpURL(packerName, packerVersion, os, arch)
		oopsBuilder          = oops.
					Code("new_packer_failed")

		absPathToPackerExecutable string
		absPathToPackerRunDir     string
		absPathToPackerZip        string
		packerZipContentLength    int64
	)

	if absPathToPackerExecutable, err = filepath.Abs(filepath.Join(dependenciesDirName, packerDirName, packerExecutableName)); err != nil {
		err = oopsBuilder.
			With("dependenciesDirName", dependenciesDirName).
			With("packerDirName", packerDirName).
			Wrapf(err, "failed to create absolute path to: %s", packerExecutableName)
		return
	}

	if absPathToPackerRunDir, err = filepath.Abs(packerDirName); err != nil {
		err = oopsBuilder.
			With("packerDirName", packerDirName).
			Wrapf(err, "failed to create absolute path to run dir")
		return
	}

	if absPathToPackerZip, err = filepath.Abs(filepath.Join(dependenciesDirName, packerDirName, packerZipName)); err != nil {
		err = oopsBuilder.
			With("dependenciesDirName", dependenciesDirName).
			With("packerDirName", packerDirName).
			Wrapf(err, "failed to create absolute path to: %s", packerZipName)
		return
	}

	if packerZipContentLength, err = utils.GetContentLength(packerUrl); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", packerUrl)
		return
	}

	packer = &Packer{
		AbsPathToExecutable: absPathToPackerExecutable,
		AbsPathToRunDir:     absPathToPackerRunDir,
		Zip: &download.Zip{
			Name:          packerZipName,
			AbsPath:       absPathToPackerZip,
			URL:           packerUrl,
			ContentLength: packerZipContentLength,
		},
	}

	return
}

func (p *Packer) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(p.AbsPathToExecutable)
}

func (p *Packer) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(p.AbsPathToExecutable)
}

func (p *Packer) SetPluginPath(cloudSetup cloud.CloudSetupI) (err error) {
	var (
		oopsBuilder = oops.
				Code("packer_set_plugin_path_failed").
				With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName())

		cloudType            = cloudSetup.GetCloudType()
		packerPluginPathName = tool.PACKER_PLUGIN_PATH_NAME
		absPluginPath        string
	)

	switch cloudType {
	case cloud.AWS:
		absPluginPath = filepath.Join(p.AbsPathToRunDir, cloud.AWS_NAME, dirs.PLUGINS_DIR_NAME)

		if err = os.Setenv(packerPluginPathName, absPluginPath); err != nil {
			err = oopsBuilder.
				With("absPluginPath", absPluginPath).
				Wrapf(err, "Error occurred while setting %s environment variable", packerPluginPathName)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloudSetup.GetCloudName())
		return
	}

	return
}

func (p *Packer) UnsetPluginPath() (err error) {
	var (
		oopsBuilder = oops.
				Code("packer_unset_plugin_path_failed")
		packerPluginPathName = tool.PACKER_PLUGIN_PATH_NAME
	)

	if err = os.Unsetenv(packerPluginPathName); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while unsetting %s environment variable", packerPluginPathName)
		return
	}

	return
}

func (p *Packer) Init() (err error) {
	var (
		cmd         = exec.Command(p.AbsPathToExecutable, "init", "-upgrade", ".")
		oopsBuilder = oops.
				Code("packer_init_failed")

		cmdErr error
	)

	if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
		err = oopsBuilder.
			Wrapf(cmdErr, "Error occured while initializing packer")
		return
	}

	return
}

func (p *Packer) Build() (err error) {
	var (
		cmd         = exec.Command(p.AbsPathToExecutable, "build", ".")
		oopsBuilder = oops.
				Code("packer_build_failed")
		cmdErr error
	)

	if cmdErr = utils.RunCmdAndStream(cmd); cmdErr != nil {
		err = oopsBuilder.
			Wrapf(cmdErr, "Error occured while building packer")
		return
	}

	return
}
