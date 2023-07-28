package instances

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Packer struct {
	ID                  binaries.Tool
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *Zip
}

const (
	PLUGINS_DIR_NAME   = "plugins"
	PACKER_PLUGIN_PATH = "PACKER_PLUGIN_PATH"
)

func NewPacker() (packer *Packer, err error) {
	const (
		PACKER  = "packer"
		VERSION = "1.9.1"
	)

	var (
		executableName = fmt.Sprintf("%s.exe", PACKER)
		zipName        = fmt.Sprintf("%s.zip", PACKER)
		os, arch       = utils.GetCurrentHostSpecs()
		url            = utils.CreateHashicorpURL(PACKER, VERSION, os, arch)
		oopsBuilder    = oops.
				Code("new_packer_failed")

		absPathToExecutable string
		absPathToRunDir     string
		zipPath             string
		contentLength       int64
	)

	if absPathToExecutable, err = filepath.Abs(filepath.Join(binaries.DEPENDENCIES_DIR_NAME, PACKER, executableName)); err != nil {
		err = oopsBuilder.
			With("DEPENDENCIES_DIR_NAME", binaries.DEPENDENCIES_DIR_NAME).
			With("PACKER", PACKER).
			Wrapf(err, "failed to create absolute path to: %s", executableName)
		return
	}

	if absPathToRunDir, err = filepath.Abs(PACKER); err != nil {
		err = oopsBuilder.
			With("PACKER", PACKER).
			Wrapf(err, "failed to create absolute path to run dir")
		return
	}

	if zipPath, err = filepath.Abs(filepath.Join(binaries.DEPENDENCIES_DIR_NAME, PACKER, zipName)); err != nil {
		err = oopsBuilder.
			With("DEPENDENCIES_DIR_NAME", binaries.DEPENDENCIES_DIR_NAME).
			With("PACKER", PACKER).
			Wrapf(err, "failed to create absolute path to: %s", zipName)
		return
	}

	if contentLength, err = utils.GetContentLength(url); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	packer = &Packer{
		ID:                  binaries.PackerID,
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

func (p *Packer) IsInstalled() (isInstalled bool) {
	return utils.FilePresent(p.AbsPathToExecutable)
}

func (p *Packer) IsNotInstalled() (isNotInstalled bool) {
	return utils.FileNotPresent(p.AbsPathToExecutable)
}

func (p *Packer) SetPluginPath(cloud binaries.Cloud) (err error) {
	var (
		oopsBuilder = oops.
				Code("packer_set_plugin_path_failed").
				With("cloud", cloud)

		absPluginPath string
	)

	switch cloud {
	case binaries.AWS:
		absPluginPath = filepath.Join(p.AbsPathToRunDir, binaries.AWS_SUBDIR_NAME, PLUGINS_DIR_NAME)

		if err = os.Setenv(PACKER_PLUGIN_PATH, absPluginPath); err != nil {
			err = oopsBuilder.
				With("absPluginPath", absPluginPath).
				Wrapf(err, "Error occurred while setting %s environment variable", PACKER_PLUGIN_PATH)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}

func (p *Packer) UnsetPluginPath(cloud binaries.Cloud) (err error) {
	var (
		oopsBuilder = oops.
			Code("packer_unset_plugin_path_failed").
			With("cloud", cloud)
	)

	switch cloud {
	case binaries.AWS:
		if err = os.Unsetenv(PACKER_PLUGIN_PATH); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while unsetting %s environment variable", PACKER_PLUGIN_PATH)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
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
