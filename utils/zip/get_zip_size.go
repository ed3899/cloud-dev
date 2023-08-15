package zip

import (
	"os"

	"github.com/samber/oops"
)

func GetZipSizeWith(
	filepathIsAbs func(string) bool,
	osOpen func(string) (*os.File, error),
) GetZipSize {
	oopsBuilder := oops.
		Code("GetZipSizeWith").
		In("utils").
		In("zip")

	getZipSize := func(absPathToZip string) (int64, error) {
		if !filepathIsAbs(absPathToZip) {
			err := oopsBuilder.
				Errorf("path to zip is not absolute: %s", absPathToZip)
			return -1, err
		}

		zipFile, err := osOpen(absPathToZip)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to open zip file: %s", absPathToZip)
			return -1, err
		}
		defer zipFile.Close()

		zipInfo, err := zipFile.Stat()
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to get zip file info: %v", absPathToZip)
			return -1, err
		}

		return zipInfo.Size(), nil
	}

	return getZipSize
}

type GetZipSize func(absPathToZip string) (size int64, err error)
