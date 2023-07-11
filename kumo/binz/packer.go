package binz

import (
	"fmt"
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
	ExecutablePath string
}

// Initializes Packer and returns the path to the Packer HCL file if successful.
func (p *Packer) init() (err error) {
	phclfp, err := utils.GetPackerHclFilePath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL file path")
		return err
	}

	cmd := exec.Command(p.ExecutablePath, "init", phclfp)
	_, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer")
		return err
	}

	return nil
}

func (p *Packer) buildAMI_OnAWS() (err error) {
	err = p.init()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer with AWS config")
		return err
	}

	// Create Packer vars files
	generalPackerVarsPath, err := templates.CraftGeneralPackerVarsFile()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing general Packer vars file")
		return err
	}

	awsPackerVarsPath, err := templates.CraftAWSPackerVarsFile()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing AWS Packer vars file")
		return err
	}

	// Change directory to Packer HCL directory
	initialLocation, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return err
	}

	runLocation, err := utils.GetPackerHclDirPath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL directory path")
		return err
	}

	err = os.Chdir(runLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to Packer HCL directory")
		return err
	}

	gVarsFileFlag := fmt.Sprintf("-var-file=%s", generalPackerVarsPath)
	awsVarsFileFlag := fmt.Sprintf("-var-file=%s", awsPackerVarsPath)

	cmd := exec.Command(p.ExecutablePath, "validate", gVarsFileFlag, awsVarsFileFlag, ".")
	cmdErr := utils.AttachToProcessStdAll(cmd)
	if cmdErr != nil {
		cmdErr = errors.Wrap(cmdErr, "Error occurred while building Packer AMI")
		err = os.Chdir(initialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change back to initial location
	err = os.Chdir(initialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (p *Packer) Build(cloud string) {
	switch cloud {
	case "aws":
		err := p.buildAMI_OnAWS()
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
	ep := filepath.Join(bins.Packer.Dependency.ExtractionPath, "packer.exe")

	// Validate existence
	if utils.FileNotPresent(ep) {
		err = errors.New("Packer executable not found")
		return nil, err
	}

	packer = &Packer{
		ExecutablePath: ep,
	}

	return packer, nil
}
