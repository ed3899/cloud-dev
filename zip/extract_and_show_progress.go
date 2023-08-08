package zip

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	utils_zip "github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
)

func ExtractAndShowProgress(
	download *Download,
	multiProgressBar interfaces.MpbV8MultiprogressBar,
	getZipSize utils_zip.GetZipSizeF,
	unzip utils_zip.UnzipF,
) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		doneChan           = make(chan bool, 1)
		oopsBuilder        = oops.
					Code("ExtractAndShowProgress").
					With("download.Name", download.Name).
					With("download.AbsPath", download.AbsPath).
					With("multiProgressBar", multiProgressBar)

		extractedBytes int
		done           bool
		zipSize        int64
	)

	if zipSize, err = getZipSize(download.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", download.AbsPath)
		return
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		download.SetExtractionBar(multiProgressBar, zipSize)

		if err = unzip(download.AbsPath, filepath.Dir(download.AbsPath), extractedBytesChan); err != nil {
			err = oopsBuilder.
				With("absPath", download.AbsPath).
				With("extractedBytesChan", extractedBytesChan).
				Wrapf(err, "failed to unzip: %v", download.AbsPath)
			errChan <- err
			return
		}

		doneChan <- true
	}(zipSize)

OuterLoop:
	for {
		select {
		case extractedBytes = <-extractedBytesChan:
			if extractedBytes > 0 {
				download.IncrementExtractionBar(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", download.Name)
				return
			}

		case done = <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return
}
