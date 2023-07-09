package utils

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func DraftDependency(name string, s Specs) (*Dependency, error) {
	destinationZipPath, err := CraftSingleZipPath(name, s)
	if err != nil {
		msg := fmt.Sprintf("failed to get zip path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	destinationExtractionPath, err := CraftSingleExtractionPath(name)
	if err != nil {
		msg := fmt.Sprintf("failed to get extraction path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	url, err := GetUrlForDep(name, s)
	if err != nil {
		msg := fmt.Sprintf("failed to get url for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	contentLength, err := GetContentLength(url)
	if err != nil {
		msg := fmt.Sprintf("failed to get content length for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	if FileNotPresent(destinationZipPath) {
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
