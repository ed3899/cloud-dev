package zip

import (
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

func GetZipSize(
	absPathToZip string,
) (int64, error) {
	oopsBuilder := oops.
		Code("GetZipSizeWith").
		In("utils").
		In("zip").
		With("absPathToZip", absPathToZip)

	if !filepath.IsAbs(absPathToZip) {
		err := oopsBuilder.
			Errorf("path to zip is not absolute: %s", absPathToZip)
		return -1, err
	}

	zipFile, err := os.Open(absPathToZip)
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
