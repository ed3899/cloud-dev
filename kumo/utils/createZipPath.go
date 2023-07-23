package utils

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func CreateZipPath(pathToZip []string) (zipPath string, err error) {
	zipPath, err = filepath.Abs(filepath.Join(pathToZip...))
	if err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %v", pathToZip)
		return
	}
	return
}
