package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

// type Terraform2I interface {
// 	Up() (err error)
// 	Down() (err error)
// }

// type Terraform2 struct {
// 	Path string
// 	Zip  *Zip
// }

// func (t *Terraform2) Up() (err error) {
// 	return
// }

// func (t *Terraform2) Down() (err error) {
// 	return
// }

type Packer2I interface {
	Build() (err error)
}

type Packer2 struct {
	ID                  Tool
	Name                string
	AbsPathToExecutable string
	Cloud               string
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
		Cloud:               cloud,
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
	// Get initial working directory
	initialLocation, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return
	}
	defer os.Chdir(initialLocation)

	//	Change working directory to where the packer files are depending on the cloud
	runLocation, err := filepath.Abs(filepath.Dir(p.AbsPathToExecutable))

	// Set PACKER_PLUGIN_PATH environment variable
	err = os.Setenv("PACKER_PLUGIN_PATH", p.AbsPathToPluginsDir)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while setting PACKER_PLUGIN_PATH environment variable")
		return
	}
	defer os.Unsetenv("PACKER_PLUGIN_PATH")

	// Initialize
	cmd := exec.Command(p.AbsPathToExecutable, "init", ".")
	if cmdErr := utils.AttachCliToProcess(cmd); cmdErr != nil {
		err = errors.Wrapf(cmdErr, "Error occured while initializing %s for %s", p.Name, p.Cloud)
		return
	}
	return
}

func (p *Packer2) Build() (err error) {

	return
}

func InitializeBinary() {

}
