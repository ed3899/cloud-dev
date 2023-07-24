package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type Packer2I interface {
	Build() (err error)
}

type Packer2 struct {
	ID                  Tool
	AbsPathToExecutable string
	AbsPathToCloudDir   string
	AbsPathToPluginsDir string
	Zip                 *Zip
}

func NewPacker(cloud string) (packer *Packer2, err error) {
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
		err = errors.Wrapf(err, "failed to create binary path to: %s", executableName)
		return
	}
	absPathToCloudDir, err := filepath.Abs(filepath.Join(name, cloud))
	if err != nil {
		err = errors.Wrapf(err, "failed to create path to: %s", cloud)
		return
	}
	absPathToPluginsDir := filepath.Join(absPathToCloudDir, "plugins")
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

	packer = &Packer2{
		ID:                  PackerID,
		AbsPathToExecutable: absPathToExecutable,
		AbsPathToPluginsDir: absPathToPluginsDir,
		Zip: &Zip{
			Name:          zipName,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}
	return
}

func (p *Packer2) Init() (err error) {
	// Set PACKER_PLUGIN_PATH environment variable
	const (
		PACKER_PLUGIN_PATH = "PACKER_PLUGIN_PATH"
	)
	err = os.Setenv(PACKER_PLUGIN_PATH, p.AbsPathToPluginsDir)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while setting %s environment variable", PACKER_PLUGIN_PATH)
		return
	}
	defer os.Unsetenv(PACKER_PLUGIN_PATH)

	// Initialize
	cmd := exec.Command(p.AbsPathToExecutable, "init", ".")
	if cmdErr := utils.AttachCliToProcess(cmd); cmdErr != nil {
		err = errors.Wrapf(cmdErr, "Error occured while initializing packer")
		return
	}
	return
}

func (p *Packer2) Build() (err error) {
	// Build
	cmd := exec.Command(p.AbsPathToExecutable, "build", ".")
	if cmdErr := utils.AttachCliToProcess(cmd); cmdErr != nil {
		err = errors.Wrapf(cmdErr, "Error occured while building packer")
		return
	}
	return
}
