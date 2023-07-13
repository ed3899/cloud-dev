package binz

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	templates_terraform "github.com/ed3899/kumo/templates/terraform"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type TerraformI interface {
	Destroy()
	Up()
}

type Terraform struct {
	InitialLocation    string
	RunLocationAbsPath string
	ExecutableAbsPath  string
}

func (t *Terraform) init(cloud string) (err error) {
	// Get and set initial location
	initialLocation, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return err
	}
	t.InitialLocation = initialLocation

	// Get and set run location
	rlap, err := utils.CraftAbsolutePath("terraform", cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to run location for cloud '%s'", cloud)
		return err
	}
	t.RunLocationAbsPath = rlap

	// Change directory to run location
	err = os.Chdir(t.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while changing directory to run location '%s'", t.RunLocationAbsPath)
		return err
	}

	// Run cmd
	cmd := exec.Command(t.ExecutableAbsPath, "init")
	cmdErr := utils.AttachCliToProcess(cmd)

	if cmdErr != nil {
		cmdErr = errors.Wrapf(cmdErr, "Error occured while initializing terraform for %v", cloud)
		// Change directory to initial location in case of error
		err = os.Chdir(t.InitialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change directory to initial location
	err = os.Chdir(t.InitialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (t *Terraform) deployToCloud(cloud string) (err error) {
	// Initialize terraform
	err = t.init(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while initializing Terraform for '%s'", cloud)
		return err
	}

	// Craft Terraform Vars file
	_, err = templates_terraform.CraftCloudTerraformTfVarsFile(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting Terraform Vars file for cloud '%s'", cloud)
		return err
	}

	// Change the directory to terraform run location
	err = os.Chdir(t.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while changing directory to Terraform directory for '%s'", cloud)
		return err
	}

	// Run cmd
	cmd := exec.Command(t.ExecutableAbsPath, "version")
	// Attach to process
	cmdErr := utils.AttachCliToProcess(cmd)
	if cmdErr != nil {
		cmdErr = errors.Wrapf(cmdErr, "Error occurred while running Terraform cmd for '%s'", cloud)
		// Change directory to initial location in case of error
		err = os.Chdir(t.InitialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			return totalError
		}
		return cmdErr
	}

	// Change back to initial location
	err = os.Chdir(t.InitialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		return err
	}

	return nil
}

func (t *Terraform) Up(cloud string) {
	err := t.deployToCloud(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while deploying to cloud '%s'", cloud)
		log.Fatal(err)
	}
}

func (t *Terraform) Destroy() {

}

func GetTerraformInstance(bins *download.Binaries) (terraform *Terraform, err error) {
	// Create the absolute path to the executable
	ep := filepath.Join(bins.Terraform.Dependency.ExtractionPath, "terraform.exe")

	// Validate existence
	if utils.FileNotPresent(ep) {
		err = errors.New("Terraform executable not found")
		return nil, err
	}

	terraform = &Terraform{
		ExecutableAbsPath: ep,
	}

	return terraform, nil
}
