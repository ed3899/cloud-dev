package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

// Crafts the absolute path to a directory where the dependency will be extracted.
//
// - {CWD}/deps/{name}
func CraftSingleExtractionPath(dirName string) (string, error) {
	cwd, err := GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}
	depsDir := GetDepsDir()

	depath := filepath.Join(cwd, depsDir, dirName)

	return depath, nil
}
