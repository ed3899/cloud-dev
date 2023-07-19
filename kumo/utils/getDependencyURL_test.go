package utils

import (
	"net/url"
	"testing"

	"github.com/ed3899/kumo/host"
)

func TestGetDependencyURL(t *testing.T) {
	// Define test cases
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

	isValidUrl := func(input string) bool {
		_, err := url.ParseRequestURI(input)
		return err == nil
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Get the dependency URL
		url, err := GetDependencyURL(tc.name, tc.specs)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Check if the URL matches the expected value
		if url != tc.expectedURL {
			t.Errorf("Expected URL: %s, Got: %s", tc.expectedURL, url)
		}

		// Check if the URL is valid
		if !isValidUrl(url) {
			t.Errorf("Expected URL to be valid, but it is not: %s", url)
		}
	}
}
