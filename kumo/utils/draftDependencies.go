package utils

import (
	"log"

	"github.com/pkg/errors"
)

func DraftDependencies(neededBinaries []string, s Specs) (*Dependencies, error) {
	missingDependencies := []*Dependency{}

	for _, v := range neededBinaries {
		msdpc, err := DraftDependency(v, s)
		if err != nil {
			error := errors.Wrap(err, "failed to draft dependency")
			log.Printf("there was an error drafting the dependency: %v", err)
			return nil, error
		}
		if msdpc == nil {
			continue
		}

		missingDependencies = append(missingDependencies, msdpc)
	}

	return &missingDependencies, nil
}
