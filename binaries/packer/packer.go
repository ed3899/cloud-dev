package packer

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/cloud"
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

func NewInstance(config binaries.ConfigI) (instance *Instance, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_failed")
		dependenciesDirName    = config.GetDependenciesDirName()
		packerName             = config.GetToolName()
		packerDirName          = packerName
		packerVersion          = config.GetToolVersion()
		packerExecutableName   = config.GetToolExecutableName()
		packerZipName          = config.GetToolZipName()
		currentOs, currentArch = config.GetCurrentHostSpecs()
		packerUrl              = config.CreateHashicorpURL(packerName, packerVersion, currentOs, currentArch)

		absPathToKumoExecutable    string
		absPathToKumoExecutableDir string
		absPathToPackerExecutable  string
		absPathToPackerRunDir      string
		absPathToPackerZip         string
		packerZipContentLength     int64
	)

	if absPathToKumoExecutable, err = config.GetKumoExecutableAbsPath(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to kumo executable")
		return
	}

	absPathToKumoExecutableDir = filepath.Dir(absPathToKumoExecutable)
	absPathToPackerExecutable = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, packerDirName, packerExecutableName)
	absPathToPackerRunDir = filepath.Join(absPathToKumoExecutableDir, packerDirName)
	absPathToPackerZip = filepath.Join(absPathToKumoExecutableDir, dependenciesDirName, packerDirName, packerZipName)

	if packerZipContentLength, err = config.GetContentLength(packerUrl); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", packerUrl)
		return
	}

	instance = &Instance{
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

func (i *Instance) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(i.AbsPathToExecutable)
}

func (i *Instance) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(i.AbsPathToExecutable)
}

func (i *Instance) SetPluginPath(cloudSetup cloud.ConfigI) (err error) {
	var (
		oopsBuilder = oops.
				Code("packer_set_plugin_path_failed").
				With("cloudSetup.GetCloudName()", cloudSetup.GetCloudName())
		cloudType            = cloudSetup.GetCloudType()
		packerPluginPathName = tool.PACKER_PLUGIN_PATH_NAME

		absPluginPath string
	)

	switch cloudType {
	case cloud.AWS:
		absPluginPath = filepath.Join(i.AbsPathToRunDir, cloud.AWS_NAME, dirs.PLUGINS_DIR_NAME)

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

func (i *Instance) UnsetPluginPath() (err error) {
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

func (i *Instance) Init() (err error) {
	var (
		cmd         = exec.Command(i.AbsPathToExecutable, "init", "-upgrade", ".")
		oopsBuilder = oops.
				Code("packer_init_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured while running and streaming packer init command")
		return
	}

	return
}

func (i *Instance) Build() (err error) {
	var (
		cmd         = exec.Command(i.AbsPathToExecutable, "build", ".")
		oopsBuilder = oops.
				Code("packer_build_failed")
	)

	if err = utils.RunCmdAndStream(cmd); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occured running and streaming packer build command")
		return
	}

	return
}
