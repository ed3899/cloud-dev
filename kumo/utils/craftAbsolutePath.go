package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Creates an absolute path with the given path(s). Cross os compatible.
// Takes the current working directory from which kumo.exe is executed as the root.
// It does not verify for existance of the resulting path.
//
// Example:
//
//	("packer", "packer.exe") -> "C:\\Users\\user\\kumo\\packer\\packer.exe"
//	("terraform", "terraform.exe") -> "/home/user/kumo/terraform/terraform.exe"
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