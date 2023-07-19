package draft

import (
	"fmt"

	"github.com/pkg/errors"
)

type Dependencies = []*Dependency

// Return a list of dependencies that are not present on the host and need to be downloaded.
func CraftDependencies() (tobeDownloaded *Dependencies, err error) {
	p, err := CraftHashicorpDependency("packer")
	if err != nil {
		msg := fmt.Sprintf("failed to draft dependency: %v", p)
		error := errors.Wrap(err, msg)
		return nil, error
	}

	t, err := CraftHashicorpDependency("terraform")
	if err != nil {
		msg := fmt.Sprintf("failed to draft dependency: %v", t)
		error := errors.Wrap(err, msg)
		return nil, error
	}

	deps := &Dependencies{p, t}
	tobeDownloaded = &Dependencies{}

	for _, d := range *deps {
		if d.Present {
			continue
		}

		*tobeDownloaded = append(*tobeDownloaded, d)
	}

	return tobeDownloaded, nil
}
