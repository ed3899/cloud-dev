package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

func GetUrlForDep(name string, s Specs) (string, error) {
	switch name {
	case "packer":
		url := fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH)
		return url, nil
	case "pulumi":
		var arch string
		switch s.ARCH {
		case "amd64":
			arch = "x64"
		}

		url := fmt.Sprintf("https://get.pulumi.com/releases/sdk/pulumi-v3.74.0-%s-%s.zip", s.OS, arch)
		return url, nil
	default:
		err := errors.New("no url found for dependency")
		return "", err
	}
}
