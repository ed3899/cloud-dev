package binz

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/ed3899/kumo/config"
)

type PackerI interface {
	Build()
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) Build(ke *config.KumoEnvironment ) {
	//TODO Add env vars to command

	log.Printf("Building AMI with Packer...")
	cmd := exec.Command(p.ExecutablePath, "build")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func GetPackerInstance(bins *utils.Binaries) (packer *Packer, err error) {
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