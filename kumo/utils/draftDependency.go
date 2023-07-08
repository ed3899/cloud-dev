package utils

import (
	"log"
	"github.com/pkg/errors"
)

func DraftDependency(name string, s Specs) (*Dependency, error) {
	destinationZipPath, err := GetDependencyZipPath(name, s)
	if err != nil {
		err = errors.Wrap(err, "failed to get zip path for dependency")
		return nil, err
	}

	destinationExtractionPath, err := GetDependencyExtractionPath(name)
	if err != nil {
		err = errors.Wrap(err, "failed to get extraction path for dependency")
		return nil, err
	}

	url, err := GetUrlForDep(name, s)
	if err != nil {
		err = errors.Wrap(err, "failed to get url for dependency")
		log.Printf("there was an error getting the url for the dependency: %v", err)
		return nil, err
	}

	contentLength, err := GetContentLength(url)
	if err != nil {
		err = errors.Wrap(err, "failed to get url for dependency")
		log.Printf("there was an error getting the url for the dependency: %v", err)
		return nil, err
	}

	if DependencyNotPresent(destinationZipPath) {
		log.Printf("%s not present", name)
		return &Dependency{
			Name:           name,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			ContentLength:  contentLength,
		}, nil
	}

	return nil, nil
}
