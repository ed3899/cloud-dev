package zip

import (
	"os"
	"path/filepath"

	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndExtractHof(zip Zip) (err error) {
	var (
		oopsBuilder = oops.
				Code("DownloadAndExtractHof")
		absPathToZipDir = filepath.Dir(zip.AbsPath)
		progress        = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
	)

	if err = os.RemoveAll(absPathToZipDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return
	}

	if err = DownloadAndShowProgress(zip, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while downloading %s", zip.Name)
		return
	}

	if err = ExtractAndShowProgress(zip, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while extracting %s", zip.Name)
		return
	}
}
