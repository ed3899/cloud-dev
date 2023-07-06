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

	// Create destination paths
	destinationZipPath := filepath.Join(dir, "deps", fmt.Sprintf("packer_%s_%s.zip", s.OS, s.ARCH))
	destinationExtractionPath := filepath.Join(dir, "packer")
	url := fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH)

	if FileExists(destinationExtractionPath) {
		log.Printf("File '%s' already exists", destinationExtractionPath)
		return &Dependency{
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			Present:        true,
		}
	}

	return &Dependency{
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		Present:        false,
	}
}
