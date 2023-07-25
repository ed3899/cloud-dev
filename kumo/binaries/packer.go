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
	PLUGINS_DIR_NAME     = "plugins"
	PACKER_PLUGIN_PATH = "PACKER_PLUGIN_PATH"
)

var (
	initialLocation string
)

func NewPacker() (packer *Packer, err error) {
	const (
		name    = "packer"
		version = "1.9.1"
	)

	var (
		executableName = fmt.Sprintf("%s.exe", name)
		zipName        = fmt.Sprintf("%s.zip", name)
		os, arch       = utils.GetCurrentHostSpecs()
		url            = utils.CreateHashicorpURL(name, version, os, arch)
		depDirName     = utils.GetDependenciesDirName()
	)

	absPathToExecutable, err := filepath.Abs(filepath.Join(depDirName, name, executableName))
	if err != nil {
		err = errors.Wrapf(err, "failed to create executable path to: %s", executableName)
		return
	}

	absPathToRunDir, err := filepath.Abs(name)
	if err != nil {
		err = errors.Wrapf(err, "failed to create run path to: %s", name)
		return
	}

	zipPath, err := filepath.Abs(filepath.Join(depDirName, name, zipName))
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
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}
	return
}

func (p *Packer) Init(cloud Cloud) (err error) {
	var (
		cmd    = exec.Command(p.AbsPathToExecutable, "init", ".")
		cmdErr error
	)

	// Store current working directory
	if initialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch cloud {
	case AWS:
		// Change directory to where packer will be run
		absPathToRunLocation := filepath.Join(p.AbsPathToRunDir, AWS_SUBDIR_NAME)
		if err = os.Chdir(absPathToRunLocation); err != nil {
			err = errors.Wrapf(err, "Error occurred while changing directory to %s", absPathToRunLocation)
			return
		}
		defer os.Chdir(initialLocation)

		// Set PACKER_PLUGIN_PATH environment variable
		err = os.Setenv(PACKER_PLUGIN_PATH, filepath.Join(absPathToRunLocation, PLUGINS_DIR_NAME))
		if err != nil {
			err = errors.Wrapf(err, "Error occurred while setting %s environment variable", PACKER_PLUGIN_PATH)
			return
		}
		defer os.Unsetenv(PACKER_PLUGIN_PATH)

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

	// Store current working directory
	if initialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch cloud {
	case AWS:
		// Change directory to where packer will be run
		absPathToRunLocation := filepath.Join(p.AbsPathToRunDir, AWS_SUBDIR_NAME)
		if err := os.Chdir(absPathToRunLocation); err != nil {
			err = errors.Wrapf(err, "Error occurred while changing directory to %s", absPathToRunLocation)
			return err
		}
		defer os.Chdir(initialLocation)

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
