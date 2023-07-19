package binz

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	templates_packer "github.com/ed3899/kumo/templates/packer"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerI interface {
	GetPackerInstance(*download.Binaries) (*Packer, error)
	Build(string)
}

type Packer struct {
	InitialLocation    string
	RunLocationAbsPath string
	ExecutableAbsPath  string
}

// Initializes Packer
func (p *Packer) init(cloud string) (err error) {
	// Create the absolute path to the plugin directory
	pdap, err := filepath.Abs(filepath.Join("packer", cloud, "plugins"))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to Packer plugins")
		return err
	}

	// Set PACKER_PLUGIN_PATH environment variable
	err = os.Setenv("PACKER_PLUGIN_PATH", pdap)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while setting PACKER_PLUGIN_PATH environment variable")
		log.Fatal(err)
	}

	// Get initial location
	initialLocation, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return err
	}
	p.InitialLocation = initialLocation

	// Get run location
	rlap, err := filepath.Abs(filepath.Join("packer", cloud))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to run location")
		return err
	}
	p.RunLocationAbsPath = rlap

	// Change directory to run location
	err = os.Chdir(p.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to run location")
		return err
	}

	// Run command
	cmd := exec.Command(p.ExecutableAbsPath, "init", ".")
	cmdErr := utils.AttachCliToProcess(cmd)

	if cmdErr != nil {
		cmdErr = errors.Wrapf(cmdErr, "Error occured while initializing packer for %v", cloud)
		// Change directory to initial location in case of error
		err = os.Chdir(p.InitialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change directory to initial location
	err = os.Chdir(p.InitialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (p *Packer) buildOnCloud(cloud string) (err error) {
	err = p.init(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while initializing Packer for '%s'", cloud)
		return err
	}

	_, err = templates_packer.CraftCloudPackerVarsFile(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting Packer vars file for '%s'", cloud)
		return err
	}

	// Change directory to Packer HCL directory
	err = os.Chdir(p.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while changing directory to Packer HCL directory for '%s'", cloud)
		return err
	}

	// Run command
	cmd := exec.Command(p.ExecutableAbsPath, "build", ".")

	// Attach to process
	cmdErr := utils.AttachCliToProcess(cmd)
	if cmdErr != nil {
		cmdErr = errors.Wrapf(cmdErr, "Error occurred while running Packer validate for '%s'", cloud)
		// Change directory to initial location in case of error
		err = os.Chdir(p.InitialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change back to initial location
	err = os.Chdir(p.InitialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (p *Packer) Build(cloud string) {
	err := p.buildOnCloud(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while building Packer AMI for %s", cloud)
		log.Fatal(err)
	}
}

func GetPackerInstance(bins *download.Binaries) (packer *Packer, err error) {
	// Create the absolute path to the executable
	eap := filepath.Join(bins.Packer.Dependency.ExtractionPath, "packer.exe")

	// Validate existence
	if utils.FileNotPresent(eap) {
		err = errors.New("Packer executable not found")
		return nil, err
	}

	packer = &Packer{
		ExecutableAbsPath: eap,
	}

	return packer, nil
}
