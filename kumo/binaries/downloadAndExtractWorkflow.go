package binaries

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndExtractWorkflow(z ZipI, extractToAbsPathDir string) (err error) {
	var (
		absPathToZipDir        = filepath.Dir(z.GetPath())
		progress                 = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
	)

	// Start with a clean slate
	if err = os.RemoveAll(extractToAbsPathDir); err != nil {
		err = errors.Wrapf(err, "Error occurred while removing %s", extractToAbsPathDir)
		return err
	}

	if err = os.RemoveAll(absPathToZipDir); err != nil {
		err = errors.Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return err
	}

	// Download
	err = DownloadAndShowProgress(z, progress)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while downloading %s", z.GetName())
		return
	}

	// Extract
	err = ExtractAndShowProgress(z, extractToAbsPathDir, progress)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while extracting %s", z.GetName())
		return
	}

	progress.Shutdown()

	// Remove zip
	err = z.Remove()
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while removing %s", z.GetName())
		return
	}

	return
}
