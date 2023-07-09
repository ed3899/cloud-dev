package binz

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func GetKumoConfig() (kc KumoConfig, err error) {
	// Get current working directory
	cwd, err := utils.GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}

	// Regex must match kumo.config.yaml or kumo.config.yml
	pattern := regexp.MustCompile(`^kumo\.config\.(yaml|yml)$`)

	// Walk the current working directory looking for the kumo config file
	// If found, set the path to the kumo config file
	// If not found, the path will remain empty string
	// This function can only prove existence of the file, not abscence!
	err = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			err = errors.Wrap(err, "failed to walk directory")
			return err
		}

		switch {
		case d.IsDir() && pattern.MatchString(d.Name()):
			return errors.New("found a directory but should be a file")
		case pattern.MatchString(d.Name()):
			log.Printf("Found kumo config file: %s", d.Name())
			kc = path
			return nil
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Check if the kumo path exists. This is somehow obvious when getting
	// no error from the above function. However, it allows us to prove
	// the abscence of the file in case of an empty string passed as the path
	if utils.FilePresent(kc) {
		return kc, nil
	}

	return kc, errors.New("kumo config file not found")
}

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
			log.Printf("Found kumo config file: %s at %s", d.Name(), path)
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
