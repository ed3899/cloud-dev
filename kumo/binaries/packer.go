package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerI interface {
	Init(cloud Cloud) (err error)
	Build() (err error)
}

type Packer struct {
	ID                  Tool
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
	)

	absPathToExecutable, err := filepath.Abs(filepath.Join(DEPENDENCIES_DIR_NAME, PACKER, executableName))
	if err != nil {
		err = errors.Wrapf(err, "failed to create executable path to: %s", executableName)
		return
	}

	absPathToRunDir, err := filepath.Abs(PACKER)
	if err != nil {
		err = errors.Wrapf(err, "failed to create run path to: %s", PACKER)
		return
	}

	zipPath, err := filepath.Abs(filepath.Join(DEPENDENCIES_DIR_NAME, PACKER, zipName))
	if err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %s", zipName)
		return
	}

	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	packer = &Packer{
		ID:                  PackerID,
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

func (p *Packer) IsInstalled() bool {
	return utils.FilePresent(p.AbsPathToExecutable)
}

func (p *Packer) IsNotInstalled() bool {
	return utils.FileNotPresent(p.AbsPathToExecutable)
}

func (p *Packer) GetAbsPathToCloudRunDir(cloud Cloud) (cloudRunDir string, err error) {
	switch cloud {
	case AWS:
		cloudRunDir = filepath.Join(p.AbsPathToRunDir, AWS_SUBDIR_NAME)
	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}
	return
}

func (p *Packer) SetPluginPath(cloud Cloud) (err error) {
	switch cloud {
	case AWS:
		if err = os.Setenv(PACKER_PLUGIN_PATH, filepath.Join(p.AbsPathToRunDir, AWS_SUBDIR_NAME, PLUGINS_DIR_NAME)); err != nil {
			err = errors.Wrapf(err, "Error occurred while setting %s environment variable", PACKER_PLUGIN_PATH)
			return
		}
	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}
	return
}

func (p *Packer) UnsetPluginPath(cloud Cloud) (err error) {
	switch cloud {
	case AWS:
		if err = os.Unsetenv(PACKER_PLUGIN_PATH); err != nil {
			err = errors.Wrapf(err, "Error occurred while unsetting %s environment variable", PACKER_PLUGIN_PATH)
			return
		}
	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
	}
	return
}

func (p *Packer) Init(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(p.AbsPathToExecutable, "init", ".")
		cmdErr error
	)

	switch cloud {
	case AWS:
		// Initialize
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrap(cmdErr, "Error occured while initializing packer")
			return
		}
		return

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
		return
	}
}

func (p *Packer) Build(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(p.AbsPathToExecutable, "build", ".")
		cmdErr error
	)

	switch cloud {
	case AWS:
		// Build
		if cmdErr = utils.AttachCliToProcess(cmd); cmdErr != nil {
			err = errors.Wrapf(cmdErr, "Error occured while building packer")
			return
		}
		return

	default:
		err = errors.Errorf("Cloud '%v' not supported", cloud)
		return
	}
}
