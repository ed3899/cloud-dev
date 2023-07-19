package utils

import (
	"os"
	"testing"
)

func TestGetLastBuiltAmiId(t *testing.T) {
	// Create a temporary packer manifest file
	tmpFile, err := os.CreateTemp("", "test_get_last_built_ami_id")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	jsonData := `{
		"builds": [
				{
						"packer_run_uuid": "abc123",
						"artifact_id": "ami:12345"
				},
				{
						"packer_run_uuid": "def456",
						"artifact_id": "ami:67890"
				}
		],
		"last_run_uuid": "def456"
	}`

	if _, err := tmpFile.WriteString(jsonData); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Run the test using the temporary packer manifest file
	amiId, err := GetLastBuiltAmiId(tmpFile.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedAmiId := "67890"
	if amiId != expectedAmiId {
		t.Errorf("Expected AMI ID: %s, Got: %s", expectedAmiId, amiId)
	}

}
