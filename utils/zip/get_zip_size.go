package zip

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

func GetZipSize(absPathToZip string) (size int64, err error) {
	var (
		oopsBuilder = oops.Code("get_zip_size_failed").
				With("absPathToZip", absPathToZip)

		zipFile *os.File
		zipInfo fs.FileInfo
	)

	if !filepath.IsAbs(absPathToZip) {
		err = oopsBuilder.
			Wrapf(err, "path to zip is not absolute: %s", absPathToZip)
		return
	}

	if zipFile, err = os.Open(absPathToZip); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to open zip file: %s", absPathToZip)
		return
	}
	defer func(zipClose *os.File) {
		if err := zipFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "failed to close zip file: %s", zipFile.Name()),
			)
		}
	}(zipFile)

	if zipInfo, err = zipFile.Stat(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip file info: %v", absPathToZip)
		return
	}

	size = zipInfo.Size()

	return
}
