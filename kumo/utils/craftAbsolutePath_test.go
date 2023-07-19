package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCraftAbsolutePath(t *testing.T) {
	// Get current working directory
	mockCwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error occurred while getting current working directory: %s", err.Error())
	}
	expectedPath := filepath.Join(mockCwd, "path1", "path2")

	// Run test
	path, err := CraftAbsolutePath("path1", "path2")
	if path != expectedPath {
		t.Errorf("Expected path: %s, Got: %s", expectedPath, path)
	}
	if err != nil {
		t.Errorf("Expected no error, Got: %s", err.Error())
	}
}
