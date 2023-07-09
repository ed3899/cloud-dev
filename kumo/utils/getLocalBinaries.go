package utils

import "github.com/pkg/errors"

// Return the local binaries
func GetLocalBinaries() (binaries *Binaries, err error) {
	// Craft extraction paths
	packerep, err := CraftSingleExtractionPath("packer")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get packer extraction path")
	}
	pulumiep, err := CraftSingleExtractionPath("pulumi")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi extraction path")
	}

	// Craft binaries
	binaries = &Binaries{
		Packer: &Binary{
			Dependency: &Dependency{
				Name:           "packer",
				ExtractionPath: packerep,
			},
		},
		Pulumi: &Binary{
			Dependency: &Dependency{
				Name:           "pulumi",
				ExtractionPath: pulumiep,
			},
		},
	}

	return binaries, nil
}