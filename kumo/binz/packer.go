package binz

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/templates"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerI interface {
	GetPackerInstance(*download.Binaries) (*Packer, error)
	Build(string)
}

type Packer struct {
	initialLocation    string
	RunLocationAbsPath string
	ExecutableAbsPath  string
}

// Initializes Packer
func (p *Packer) init(cloud string) (err error) {
	// Create the absolute path to the plugin directory
	pdaa, err := utils.CraftAbsolutePath("packer", cloud, "plugins")
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to Packer plugins")
		return err
	}

	// Set PACKER_PLUGIN_PATH environment variable
	err = os.Setenv("PACKER_PLUGIN_PATH", pdaa)
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
	p.initialLocation = initialLocation

	// Get run location
	rlap, err := utils.CraftAbsolutePath("packer", cloud)
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
	_, err = cmd.CombinedOutput()

	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer")

		// Change directory to initial location in case of error
		errx := os.Chdir(p.initialLocation)
		if errx != nil {
			errx = errors.Wrap(err, "Error occurred while changing directory to initial location")
			return errx
		}

		return err
	}

	// Change directory to initial location
	err = os.Chdir(p.initialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (p *Packer) buildAMI_OnCloud(cloud string) (err error) {
	err = p.init(cloud)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer with AWS config")
		return err
	}

	// Create Packer vars files
	_, err = templates.CraftGeneralPackerVarsFile()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing general Packer vars file")
		return err
	}

	_, err = templates.CraftAWSPackerVarsFile()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing AWS Packer vars file")
		return err
	}

	// Change directory to Packer HCL directory
	err = os.Chdir(p.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to Packer HCL directory")
		return err
	}

	// Run command
	cmd := exec.Command(p.ExecutableAbsPath, "validate", ".")

	// Attach to process
	cmdErr := utils.AttachToProcessStdAll(cmd)
	if cmdErr != nil {
		cmdErr = errors.Wrap(cmdErr, "Error occurred while building Packer AMI")
		// Change directory to initial location in case of error
		err = os.Chdir(p.initialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change back to initial location
	err = os.Chdir(p.initialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (p *Packer) Build(cloud string) {
	switch cloud {
	case "aws":
		err := p.buildAMI_OnCloud(cloud)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while building AMI on AWS")
			log.Fatal(err)
		}
	default:
		err := errors.Errorf("Cloud '%s' not supported", cloud)
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
