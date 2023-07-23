package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func CreateBinaryPath(pathToBin []string) (absBinPath string, err error) {
	absBinPath, err = filepath.Abs(filepath.Join(pathToBin...))
	if err != nil {
		err = errors.Wrapf(err, "failed to create binary path to: %v", pathToBin)
		return
	}
	return
}
