package utils

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
)

func DraftDependency(name string, s Specs) (*Dependency, error) {
	cwd, err := GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		log.Printf("there was an error getting the current directory: %v", err)
		return nil, err
	}

	depsDir := "deps"

	getDestinationZipPath := func() string {
		return filepath.Join(cwd, depsDir, fmt.Sprintf("%s_%s_%s.zip", name, s.OS, s.ARCH))
	}

	getDestinationExtractionPath := func() string {
		return filepath.Join(cwd, depsDir, name)
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

	if DependencyNotPresent(getDestinationZipPath()) {
		log.Printf("%s not present", name)
		return &Dependency{
			Name:           name,
			URL:            url,
			ExtractionPath: getDestinationExtractionPath(),
			ZipPath:        getDestinationZipPath(),
			ContentLength:  contentLength,
		}, nil
	}

	return nil, nil
}
