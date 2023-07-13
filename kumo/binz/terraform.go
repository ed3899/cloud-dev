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
	"github.com/spf13/viper"
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

func (t *Terraform) setInitialAndRunLocations(cloud string) (err error) {
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

	return nil
}

func (t *Terraform) init(cloud string) (err error) {
	// Set initial and run locations
	err = t.setInitialAndRunLocations(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while setting initial and run locations for cloud '%s'", cloud)
		return err
	}

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

func setCloudCredentials(cloud string) (err error) {
	switch cloud {
	case "aws":
		err = os.Setenv("AWS_ACCESS_KEY_ID", viper.GetString("AWS.AccessKeyId"))
		if err != nil {
			err = errors.Wrap(err, "Error occurred while setting AWS_ACCESS_KEY_ID environment variable")
			return err
		}

		err = os.Setenv("AWS_SECRET_ACCESS_KEY", viper.GetString("AWS.SecretAccessKey"))
		if err != nil {
			err = errors.Wrap(err, "Error occurred while setting AWS_SECRET_ACCESS_KEY environment variable")
			return err
		}

		return nil

	default:
		err = errors.Errorf("Cloud '%s' is not supported", cloud)
		return err
	}
}

func (t *Terraform) deployToCloud(cloud string) (err error) {
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
	cmd := exec.Command(t.ExecutableAbsPath, "apply", "-auto-approve")
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
	// Initialize terraform
	err := t.init(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while initializing Terraform for '%s'", cloud)
		log.Fatal(err)
	}

	// Set cloud credentials
	err = setCloudCredentials(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while setting cloud credentials for '%s'", cloud)
		log.Fatal(err)
	}

	// Deploy to cloud
	err = t.deployToCloud(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while deploying to cloud '%s'", cloud)
		log.Fatal(err)
	}
}

func (t *Terraform) Destroy(cloud string) {
	// Set cloud credentials
	err := setCloudCredentials(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while setting cloud credentials for '%s'", cloud)
		log.Fatal(err)
	}

	// Set initial and run locations
	err = t.setInitialAndRunLocations(cloud)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while setting initial and run locations for cloud '%s'", cloud)
		log.Fatal(err)
	}

	// Change the directory to terraform run location
	err = os.Chdir(t.RunLocationAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while changing directory to Terraform directory for '%s'", cloud)
		log.Fatal(err)
	}

	// Run cmd
	cmd := exec.Command(t.ExecutableAbsPath, "destroy", "-auto-approve")
	// Attach to process
	cmdErr := utils.AttachCliToProcess(cmd)
	if cmdErr != nil {
		cmdErr = errors.Wrapf(cmdErr, "Error occurred while running Terraform cmd for '%s'", cloud)
		// Change directory to initial location in case of error
		err = os.Chdir(t.InitialLocation)
		if err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory to initial location")
			totalError := errors.Wrap(cmdErr, err.Error())
			log.Fatal(totalError)
		}
		log.Fatal(cmdErr)
	}

	// Change back to initial location
	err = os.Chdir(t.InitialLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to initial location")
		log.Fatal(err)
	}
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
