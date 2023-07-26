package download

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndExtract(z binaries.ZipI, extractToAbsPathDir string) (err error) {
	var (
		absPathToZipDir = filepath.Dir(z.GetPath())
		progress        = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
	)

	// Start with a clean slate
	if err = os.RemoveAll(extractToAbsPathDir); err != nil {
		return errors.Wrapf(err, "Error occurred while removing %s", extractToAbsPathDir)
	}

	if err = os.RemoveAll(absPathToZipDir); err != nil {
		return errors.Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
	}

	// Download
	if err = DownloadAndShowProgress(z, progress); err != nil {
		return errors.Wrapf(err, "Error occurred while downloading %s", z.GetName())
	}

	// Extract
	if err = ExtractAndShowProgress(z, extractToAbsPathDir, progress); err != nil {
		return errors.Wrapf(err, "Error occurred while extracting %s", z.GetName())
	}

	progress.Shutdown()

	// Remove zip
	if err = z.Remove(); err != nil {
		return errors.Wrapf(err, "Error occurred while removing %s", z.GetName())
	}

	return
}
