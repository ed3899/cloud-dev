package draft

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCraftDependency(t *testing.T) {
	// Testing a basic dependency with expected props
	expectDepName := "packer"

	dependency, err := CraftHashicorpDependency("packer")
	if err != nil {
		t.Errorf("Error while crafting dependency: %v", err)
	}

	if dependency.Name != expectDepName {
		t.Errorf("Expected name: %s, but got: %s", expectDepName, dependency.Name)
	}

	if dependency.Present {
		t.Errorf("Expected dependency to be not present")
	}

	if dependency.URL == "" {
		t.Errorf("Expected URL not to be empty")
	}

	if dependency.ExtractionPath == "" {
		t.Errorf("Expected extraction path not to be empty")
	}

	if dependency.ZipPath == "" {
		t.Errorf("Expected zip path not to be empty")
	}

	if dependency.ContentLength == 0 {
		t.Errorf("Expected content length not to be empty")
	}

	// -> Create a mock dependency
	depsDir, err := filepath.Abs("deps")
	if err != nil {
		t.Fatalf("Failed to craft absolute path: %v", err)
	}
	execPath := filepath.Join(depsDir, "packer", "packer.exe")
	if err := os.MkdirAll(execPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Test an existing dependency
	dependency, err = CraftHashicorpDependency("packer")
	if err != nil {
		t.Errorf("Error while crafting dependency: %v", err)
	}

	if !dependency.Present {
		t.Errorf("Expected dependency to be present")
	}

	// -> Remove the mock dependency
	if err := os.RemoveAll(depsDir); err != nil {
		t.Fatalf("Failed to remove directory: %v", err)
	}

	// Test a non-existing dependency
	dependency, err = CraftHashicorpDependency("packer")
	if err != nil {
		t.Errorf("Error while crafting dependency: %v", err)
	}

	if dependency.Present {
		t.Errorf("Expected dependency not to be present")
	}

}
