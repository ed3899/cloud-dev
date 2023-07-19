package utils

import (
	"testing"

	"github.com/ed3899/kumo/host"
)

func TestGetDependencyURL(t *testing.T) {
	testCases := []struct {
		specs       *host.Specs
		name        string
		version     string
		expectedURL string
	}{
		{
			specs: &host.Specs{
				OS:   "windows",
				ARCH: "amd64",
			},
			name:        "packer",
			expectedURL: "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_windows_amd64.zip",
		},
		{
			specs: &host.Specs{
				OS:   "windows",
				ARCH: "386",
			},
			name:        "packer",
			expectedURL: "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_windows_386.zip",
		},
		{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "amd64",
			},
			name:        "packer",
			expectedURL: "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_darwin_amd64.zip",
		},
		{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "arm64",
			},
			name:        "packer",
			expectedURL: "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_darwin_arm64.zip",
		},
	}

	for _, tc := range testCases {
		url, err := GetDependencyURL(tc.name, tc.specs)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if url != tc.expectedURL {
			t.Errorf("Expected URL: %s, Got: %s", tc.expectedURL, url)
		}
	}
}
