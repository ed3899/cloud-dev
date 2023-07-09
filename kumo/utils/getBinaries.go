package utils

import (
	"github.com/pkg/errors"
)

func GetBinaries() (binaries *Binaries, err error) {
	// Validate host compatibility
	specs := ValidateHostCompatibility(GetHostSpecs())

	// Get missing dependencies
	neededBinaries := []string{"packer", "pulumi"}
	missingDependencies, err := DraftDependencies(neededBinaries, specs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to draft missing dependencies")
	}

	// No missing dependencies
	if len(*missingDependencies) == 0 {
		// Get already downloaded binaries
		binaries, err := GetLocalBinaries()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get local binaries")
		}
		return binaries, nil
	}

	// Download missing dependencies
	binaries, err = DownloadDependencies(missingDependencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download dependencies")
	}

	return binaries, nil
}
