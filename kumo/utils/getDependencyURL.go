package utils

import (
	"fmt"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

// Returns the latest version of Packer. Value is hardcoded for now.
func GetLatestPackerVersion() string {
	return "1.9.1"
}

// Returns the latest version of Terraform. Value is hardcoded for now.
func GetLatestTerraformVersion() string {
	return "1.5.3"
}

// Returns the URL for the dependency specified by name.
//
// Examples:
//
//	s := host.Specs{OS: "windows", ARCH: "amd64"}
//	("packer", s) -> "https://releases.hashicorp.com/packer/1.7.4/packer_1.7.4_windows_amd64.zip", nil
//	("terraform", s) -> "https://releases.hashicorp.com/terraform/1.0.0/terraform_1.0.0_windows_amd64.zip", nil
func GetDependencyURL(name string, s *host.Specs) (url string, err error) {
	switch name {
	case "packer":
		version := GetLatestPackerVersion()
		url = fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_%s_%s.zip", version, version, s.OS, s.ARCH)
		return url, nil
	case "terraform":
		version := GetLatestTerraformVersion()
		url = fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip", version, version, s.OS, s.ARCH)
		return url, nil
	default:
		err := errors.New("no url found for dependency")
		return "", err
	}
}
