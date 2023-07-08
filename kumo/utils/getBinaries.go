package utils

import (
	"github.com/pkg/errors"
)

func GetBinaries() (*Binaries, error) {
	// Validate host compatibility
	specs := ValidateHostCompatibility(GetHostSpecs())

	// Get missing dependencies
	neededBinaries := []string{"packer", "pulumi"}
	missingDependencies, err := DraftDependencies(neededBinaries, specs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to draft missing dependencies")
	}

	// Download missing dependencies
	binaries, err := DownloadDependencies(missingDependencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download dependencies")
	}

	return binaries, nil
}
