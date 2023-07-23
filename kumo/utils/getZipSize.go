package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetZipSize(absPathToZip string) (size int64, err error) {
	if filepath.IsLocal(absPathToZip) {
		err = errors.New("path to zip is not absolute")
		return
	}

	zipfile, err := os.Open(absPathToZip)
	if err != nil {
		err = errors.Wrapf(err, "failed to open zip file: %v", absPathToZip)
		return
	}
	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		err = errors.Wrapf(err, "failed to get zip file info: %v", absPathToZip)
		return
	}

	size = info.Size()
	return
}
