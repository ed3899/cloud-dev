package draft

import (
	"fmt"
	"log"

	"github.com/ed3899/kumo/host"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func DraftTerraformDependency() (dp *Dependency, err error) {
	name := "terraform"
	s := host.GetSpecs()

	destinationZipPath, err := utils.CraftSingleZipPath(name, s)
	if err != nil {
		msg := fmt.Sprintf("failed to get zip path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	destinationExtractionPath, err := utils.CraftSingleExtractionPath(name)
	if err != nil {
		msg := fmt.Sprintf("failed to get extraction path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	url, err := utils.GetDependencyURL(name, s)
	if err != nil {
		msg := fmt.Sprintf("failed to get url for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		msg := fmt.Sprintf("failed to get content length for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	executablePath, err := utils.GetExecutablePath(name)
	if err != nil {
		msg := fmt.Sprintf("failed to get executable path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	if utils.FileNotPresent(executablePath) {
		log.Printf("%s not present", name)
		return &Dependency{
			Name:           name,
			Present:        false,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			ContentLength:  contentLength,
		}, nil
	}

	return &Dependency{
		Name:           name,
		Present:        true,
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		ContentLength:  contentLength,
	}, nil

}