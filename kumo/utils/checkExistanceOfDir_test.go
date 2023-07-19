package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDirPath := filepath.Join(cwd, "tmp")
	err = os.Mkdir(tempDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.Remove(tempDirPath)

	// Test existance
	exist, err := DirExist(tempDirPath)
	if err != nil {
		t.Errorf("Unexpected error occurred: %v", err)
	}
	if !exist {
		t.Errorf("Expected directory to exist, but it does not")
	}
}

func TestDirNotExists(t *testing.T)  {
	// Craft a path to a non-existing directory
	nonExistantDirPath := "path/to/nonexisting/dir"
	
	// Test non existance
	notExist, err := DirNotExist(nonExistantDirPath)
	exists := !notExist
	if err != nil {
		t.Errorf("Unexpected error occurred: %v", err)
	}
	if exists {
		t.Errorf("Expected directory to not exist, but it does")
	}
}
