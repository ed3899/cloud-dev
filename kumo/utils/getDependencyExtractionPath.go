package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func GetDependencyExtractionPath(name string) (string, error) {
	cwd, err := GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}
	depsDir := GetDepsDir()

	depath := filepath.Join(cwd, depsDir, name)

	return depath, nil
}
