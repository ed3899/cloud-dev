package utils

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetZipSize(absPathToZip string) (size int64, err error) {
	var (
		zipFile *os.File
		zipInfo fs.FileInfo
	)

	if filepath.IsLocal(absPathToZip) {
		err = errors.New("path to zip is not absolute")
		return
	}

	if zipFile, err = os.Open(absPathToZip); err != nil {
		err = errors.Wrapf(err, "failed to open zip file: %v", absPathToZip)
		return
	}
	defer func() {
		if errClosingZipFile := zipFile.Close(); errClosingZipFile != nil {
			err = errors.Wrapf(errClosingZipFile, "failed to close zip file: %v", absPathToZip)
		}
	}()

	if zipInfo, err = zipFile.Stat(); err != nil {
		err = errors.Wrapf(err, "failed to get zip file info: %v", absPathToZip)
		return
	}

	size = zipInfo.Size()

	return
}
