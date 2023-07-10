package binz

import (
	"log"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/host"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func GetBinaries() (binaries *download.Binaries) {
	// Get missing dependencies
	neededBinaries := []string{"packer", "pulumi"}
	missingDependencies, err := draft.DraftDependencies(neededBinaries, host.GetSpecs())
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
	binaries, err = download.DownloadDependencies(missingDependencies)
	if err != nil {
		err = errors.Wrap(err, "failed to download dependencies")
		log.Fatal(err)
	}

	return binaries
}

func GetBinaryInstances(bins *download.Binaries) (packer *Packer, pulumi *Pulumi) {
	packer, err := GetPackerInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer instance")
		log.Fatal(err)
	}

	pulumi, err = GetPulumiInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Pulumi instance")
		log.Fatal(err)
	}

	return packer, pulumi
}

// Return the local binaries
func GetLocalBinaries() (binaries *download.Binaries, err error) {
	// Craft extraction paths
	packerep, err := utils.CraftSingleExtractionPath("packer")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get packer extraction path")
	}
	pulumiep, err := utils.CraftSingleExtractionPath("pulumi")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi extraction path")
	}

	// Craft binaries
	binaries = &download.Binaries{
		Packer: &download.Binary{
			Dependency: &draft.Dependency{
				Name:           "packer",
				ExtractionPath: packerep,
			},
		},
		Pulumi: &download.Binary{
			Dependency: &draft.Dependency{
				Name:           "pulumi",
				ExtractionPath: pulumiep,
			},
		},
	}

	return binaries, nil
}