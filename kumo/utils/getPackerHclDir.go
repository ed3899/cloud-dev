package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Returns the absolute path to the packer hcl directory
func GetPackerHclDirPath() (phcldirpath string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	phcldirpath = filepath.Join(cwd, "packer")
	return phcldirpath, nil
}
