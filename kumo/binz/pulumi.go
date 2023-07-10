package binz

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PulumiI interface {
	Up()
}

type Pulumi struct {
	ExecutablePath string
}

func (p *Pulumi) deployToAWS() {
	cmd := exec.Command(p.ExecutablePath, "version")
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while building AMI")
		log.Fatal(err)
	}
}

func (p *Pulumi) Up(cloud string) {
	// Build the kumo deployment
	//
	switch cloud {
	case "aws":
		p.deployToAWS()
	default:
		err := errors.Errorf("Cloud '%s' not supported", cloud)
		log.Fatal(err)
	}
}

func (p *Pulumi) Destroy()  {

}

func GetPulumiInstance(bins *download.Binaries) (pulumi *Pulumi, err error) {
	// Create the absolute path to the executable
	ep := filepath.Join(bins.Pulumi.Dependency.ExtractionPath, "pulumi", "bin", "pulumi.exe")

	// Validate existence
	if utils.FileNotPresent(ep) {
		err = errors.New("Pulumi executable not found")
		return nil, err
	}

	pulumi = &Pulumi{
		ExecutablePath: ep,
	}

	return pulumi, nil
}
