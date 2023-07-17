package binz

import (
	"log"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func GetBinaries() (binaries *download.Binaries) {
	// Get missing dependencies
	missingDependencies, err := draft.CraftDependencies()
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

	// Remove downloaded zips
	err = download.RemoveDownloads(missingDependencies)
	if err != nil {
		err = errors.Wrapf(err, "failed to remove downloads")
		log.Print(err)
	}

	return binaries
}

func GetBinaryInstances(bins *download.Binaries) (packer *Packer, terraform *Terraform) {
	packer, err := GetPackerInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer instance")
		log.Fatal(err)
	}

	terraform, err = GetTerraformInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Terraform instance")
		log.Fatal(err)
	}

	return packer, terraform
}

// Return the local binaries
func GetLocalBinaries() (binaries *download.Binaries, err error) {
	depsDir := utils.GetDependenciesDirName()

	// Craft extraction paths
	packerep, err := utils.CraftAbsolutePath(depsDir, "packer")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get packer extraction path")
	}
	terraformep, err := utils.CraftAbsolutePath(depsDir, "terraform")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get terraform extraction path")
	}

	// Craft binaries
	binaries = &download.Binaries{
		Packer: &download.Binary{
			Dependency: &draft.Dependency{
				Name:           "packer",
				ExtractionPath: packerep,
			},
		},
		Terraform: &download.Binary{
			Dependency: &draft.Dependency{
				Name:           "terraform",
				ExtractionPath: terraformep,
			},
		},
	}

	return binaries, nil
}
