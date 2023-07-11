package binz

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerI interface {
	DraftDependency() (*draft.Dependency, error)
	GetPackerInstance(*download.Binaries) (*Packer, error)
	Build(string)
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) buildAMI_OnAWS() {
	cmd := exec.Command(p.ExecutablePath, "version")
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while building AMI")
		log.Fatal(err)
	}
}

func (p *Packer) Build(cloud string) {
	switch cloud {
	case "aws":
		p.buildAMI_OnAWS()
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
