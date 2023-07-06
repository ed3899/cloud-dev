package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetPackerUrl(s Specs) *ZipExecutableRef {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("there was an error getting the current directory: %v", err)
	}

	// Create the destination path with dir + packer + packer.exe
	destinationZipPath := filepath.Join(dir, "packer", fmt.Sprintf("packer_%s_%s.zip", s.OS, s.ARCH))

	// Return the zip executable reference
	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH),
		BinPath: destinationZipPath,
	}
}