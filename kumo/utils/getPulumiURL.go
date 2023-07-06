package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetPulumiUrl(s Specs) *ZipExecutableRef {
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
	destinationZipPath := filepath.Join(dir, "pulumi", fmt.Sprintf("packer_%s_%s.zip", s.OS, arch))

	// Return the zip executable reference
	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://get.pulumi.com/releases/sdk/pulumi-v3.74.0-%s-%s.zip", s.OS, arch),
		BinPath: destinationZipPath,
	}
}
