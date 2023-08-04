package download

import (
	"os"
	"path/filepath"

	common_zip_interfaces "github.com/ed3899/kumo/common/zip/interfaces"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func New(zip common_zip_interfaces.Zip, extractToAbsPathDir string) (err error) {
	var (
		absPathToZipDir = filepath.Dir(zip.AbsPath())
		progress        = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
		oopsBuilder     = oops.
				Code("common_download_new_failed").
				With("extractToAbsPathDir", extractToAbsPathDir).
				With("zip", zip.Name()).
				With("absPath", zip.AbsPath())
	)

	// Start with a clean slate
	if err = os.RemoveAll(extractToAbsPathDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", extractToAbsPathDir)
		return
	}

	if err = os.RemoveAll(absPathToZipDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return
	}

	// Download
	if err = DownloadAndShowProgress(zip, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while downloading %s", zip.Name())
		return
	}

	// Extract
	if err = ExtractAndShowProgress(zip, extractToAbsPathDir, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while extracting %s", zip.Name())
		return
	}

	progress.Shutdown()

	// Remove zip
	if err = zip.Remove(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", zip.Name())
		return
	}

	return
}
