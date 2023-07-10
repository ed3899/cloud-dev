package utils

import (
	"log"

	"github.com/pkg/errors"
)

func GetBinaries() (binaries *Binaries) {
	// Validate host compatibility
	specs := ValidateHostCompatibility(GetHostSpecs())

	// Get missing dependencies
	neededBinaries := []string{"packer", "pulumi"}
	missingDependencies, err := DraftDependencies(neededBinaries, specs)
	if err != nil {
		err = errors.Wrap(err, "failed to draft missing dependencies")
		log.Fatal(err)
	}

	// No missing dependencies
	if len(*missingDependencies) == 0 {
		// Get already downloaded binaries
		binaries, err := GetLocalBinaries()
		if err != nil {
			err = errors.Wrap(err, "failed to get local binaries")
			log.Fatal(err)
		}
		return binaries
	}

	// Download missing dependencies
	binaries, err = DownloadDependencies(missingDependencies)
	if err != nil {
		err = errors.Wrap(err, "failed to download dependencies")
		log.Fatal(err)
	}

	return binaries
}
