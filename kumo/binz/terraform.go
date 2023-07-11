package binz

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type TerraformI interface {
	Up()
}

type Terraform struct {
	ExecutablePath string
}



func (t *Terraform) deployToAWS() {
	cmd := exec.Command(t.ExecutablePath, "version")
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while building AMI")
		log.Fatal(err)
	}
}

func (t *Terraform) Up(cloud string) {
	// Build the kumo deployment
	//
	switch cloud {
	case "aws":
		t.deployToAWS()
	default:
		err := errors.Errorf("Cloud '%s' not supported", cloud)
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
		ExecutablePath: ep,
	}

	return terraform, nil
}
