package utils

import (
	"fmt"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

func GetDependencyURL(name string, s *host.Specs) (url string, err error) {
	switch name {
	case "packer":
		url = fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH)
		return url, nil
	case "terraform":
		url = fmt.Sprintf("https://releases.hashicorp.com/terraform/1.5.2/terraform_1.5.2_%s_%s.zip", s.OS, s.ARCH)
		return url, nil
	default:
		err := errors.New("no url found for dependency")
		return "", err
	}
}
