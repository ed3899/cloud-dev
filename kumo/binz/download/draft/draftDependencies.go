package draft

import (
	"fmt"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

type Dependencies = []*Dependency

// Return a list of dependencies that are not present on the host and need to be downloaded.
func DraftDependencies(depsList []string, s *host.Specs) (deps *Dependencies, err error) {
	deps = &Dependencies{}

	for _, d := range depsList {
		dd, err := DraftDependency(d, s)
		if err != nil {
			msg := fmt.Sprintf("failed to draft dependency: %v", d)
			error := errors.Wrap(err, msg)
			return nil, error
		}
		if dd.Present {
			continue
		}

		*deps = append(*deps, dd)
	}

	return deps, nil
}
