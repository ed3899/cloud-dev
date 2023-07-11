package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Returns the absolute path to the packer hcl file
func GetPackerHclFilePath() (phclpath string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	phclpath = filepath.Join(cwd, "packer", "ami.pkr.hcl")

	if FileNotPresent(phclpath) {
		err = errors.New("Packer HCL file not found")
		return "", err
	}

	return phclpath, nil
}
