package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Crafts the absolute path to a directory where the dependency will be extracted.
//
// - {CWD}/deps/{name}
func CraftSingleExtractionPath(dirName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}
	depsDir := GetDependenciesDirName()

	depath := filepath.Join(cwd, depsDir, dirName)

	return depath, nil
}
