package utils

import (
	"os"
	"testing"
)

func TestFilePresent(t *testing.T) {
	// Create a temporary file for testing
	existingTempFile, err := os.CreateTemp("", "test-file-present")
	if err != nil {
		t.Errorf("Error occurred while creating temp file: %s", err.Error())
	}
	// Get the file path of the temporary file
	etfp := existingTempFile.Name()
	defer os.Remove(etfp)

	// Test existance
	if !FilePresent(etfp) {
		t.Errorf("Expected file %s to exist, but it does not", etfp)
	}
}

func TestFileNotPresent(t *testing.T) {
	// Craft a path to a non-existing file
	nonExistingFilePath := "path/to/nonexisting/file"
	// Test non existance
	if !FileNotPresent(nonExistingFilePath) {
		t.Errorf("Expected file %s to be not present, but it is", nonExistingFilePath)
	}
}
