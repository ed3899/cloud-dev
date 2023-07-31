package download

import (
	"os"
	"path/filepath"

	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func Initiate(z ZipI, extractToAbsPathDir string) (err error) {
	var (
		absPathToZipDir = filepath.Dir(z.GetPath())
		progress        = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
		oopsBuilder     = oops.
				Code("initiate_failed").
				With("extractToAbsPathDir", extractToAbsPathDir).
				With("z.GetName()", z.GetName()).
				With("z.GetPath()", z.GetPath())
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
	if err = DownloadAndShowProgress(z, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while downloading %s", z.GetName())
		return
	}

	// Extract
	if err = ExtractAndShowProgress(z, extractToAbsPathDir, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while extracting %s", z.GetName())
		return
	}

	progress.Shutdown()

	// Remove zip
	if err = z.Remove(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", z.GetName())
		return
	}

	return
}