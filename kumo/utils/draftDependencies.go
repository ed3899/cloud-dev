package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

func DraftDependencies(neededBinaries []string, s Specs) (*Dependencies, error) {
	missingDependencies := []*Dependency{}

	for _, v := range neededBinaries {
		msdpc, err := DraftDependency(v, s)
		if err != nil {
			msg := fmt.Sprintf("failed to draft dependency: %v", v)
			error := errors.Wrap(err, msg)
			return nil, error
		}
		if msdpc == nil {
			continue
		}

		missingDependencies = append(missingDependencies, msdpc)
	}

	return &missingDependencies, nil
}
