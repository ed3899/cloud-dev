package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DraftPackerDependency(s Specs) *Dependency {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("there was an error getting the current directory: %v", err)
	}

	// Create fields
	name := "packer"
	destinationZipPath := filepath.Join(dir, "deps", fmt.Sprintf("%s_%s_%s.zip", name, s.OS, s.ARCH))
	destinationExtractionPath := filepath.Join(dir, name)
	url := fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH)
	contentLength := GetContentLength(url)

	// Check if the file already exists
	if FileExists(destinationExtractionPath) {
		log.Printf("File '%s' already exists", destinationExtractionPath)
		return &Dependency{
			Name:           name,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			Present:        true,
			ContentLength:  contentLength,
		}
	}

	// Return the dependency
	return &Dependency{
		Name:           name,
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		Present:        false,
		ContentLength:  contentLength,
	}
}
