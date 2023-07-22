package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func CreateBinaryPath(pathToBin []string) (absBinPath string, err error) {
	absBinPath, err = filepath.Abs(filepath.Join(pathToBin...))
	if err != nil {
		err = errors.Wrap(err, "failed to get absolute path for packer binary")
		return
	}
	return
}