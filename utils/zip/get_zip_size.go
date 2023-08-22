package zip

import (
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

// Returns the size of the zip file at the given path.
// The path must be absolute.
func GetZipSize(
	pathToZip string,
) (int64, error) {
	oopsBuilder := oops.
		Code("GetZipSize").
		In("utils").
		In("zip").
		With("pathToZip", pathToZip)

	if !filepath.IsAbs(pathToZip) {
		err := oopsBuilder.
			Errorf("path to zip is not absolute: %s", pathToZip)
		return -1, err
	}

	zipFile, err := os.Open(pathToZip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to open zip file: %s", pathToZip)
		return -1, err
	}
	defer zipFile.Close()

	zipInfo, err := zipFile.Stat()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get zip file info: %v", pathToZip)
		return -1, err
	}

	return zipInfo.Size(), nil
}
