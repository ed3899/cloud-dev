package binz

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

// Checks if the kumo config file is present in the current working directory
func KumoConfigPresent() (bool, error) {
	// Get current working directory
	cwd, err := utils.GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return false, err
	}

	// Regex must match kumo.config.yaml or kumo.config.yml
	pattern := regexp.MustCompile(`^kumo\.config\.(yaml|yml)$`)

	// Walk the current working directory looking for the kumo config file
	err = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			err = errors.Wrap(err, "failed to walk directory")
			return err
		}

		switch {
		case d.IsDir() && pattern.MatchString(d.Name()):
			log.Println(d.Name())
			return errors.New("found a directory but should be a file")
		case pattern.MatchString(d.Name()):
			log.Printf("Found kumo config file: %s", d.Name())
			return nil
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// Checks if the kumo config file is not present in the current working directory
func KumoConfigNotPresent() (bool, error) {
	present, err := KumoConfigPresent()
	if err != nil {
		err = errors.Wrap(err, "failed to check if kumo config is present")
		return false, err
	}

	return !present, nil
}
