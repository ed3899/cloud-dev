package utils

import (
	"fmt"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

func GetLastPackerVersion() string {
	return "1.9.1"
}

func GetLatestTerraformVersion() string {
	return "1.5.3"
}

func GetDependencyURL(name string, s *host.Specs) (url string, err error) {
	switch name {
	case "packer":
		version := GetLastPackerVersion()
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
