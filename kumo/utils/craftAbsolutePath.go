package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Creates an absolute path with the given path(s). Cross os compatible.
// Takes the current working directory from which main.go is executed, or if called
// from a test, the current working directory from which the test is executed.
// It does not verify for existance of the resulting path.
//
// Example:
//
//	("path1", "path2") -> "C:\\Users\\user\\kumo\\path1\\path2"
//	("templates") -> "C:\\Users\\user\\kumo\\templates"
func CraftAbsolutePath(paths ...string) (absolutePath string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	absolutePath = filepath.Join(cwd, filepath.Join(paths...))

	return absolutePath, nil
}