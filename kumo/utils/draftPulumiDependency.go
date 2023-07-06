package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DraftPulumiDependency(s Specs) *Dependency {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("there was an error getting the current directory: %v", err)
	}

	// Create the destination path
	var arch string
	switch s.ARCH {
	case "amd64":
		arch = "x64"
	}
	// Create fields
	name := "pulumi"
	destinationZipPath := filepath.Join(dir, "deps", fmt.Sprintf("%s_%s_%s.zip", name, s.OS, arch))
	destinationExtractionPath := filepath.Join(dir, name)
	url := fmt.Sprintf("https://get.pulumi.com/releases/sdk/pulumi-v3.74.0-%s-%s.zip", s.OS, arch)
	contentLength := GetContentLength(url)

	if FileExists(destinationExtractionPath) {
		log.Printf("File '%s' already exists", destinationExtractionPath)
		return &Dependency{
			Name:           name,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			Present:        true,
			ContentLength: contentLength,
		}
	}

	return &Dependency{
		Name:           name,
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		Present:        false,
		ContentLength: contentLength,
	}

}
