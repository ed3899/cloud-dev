package binz

import (
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func GetPulumiInstance(bins *utils.Binaries) (pulumi *Pulumi, err error) {
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
