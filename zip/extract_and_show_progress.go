package zip

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	utils_zip "github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
)

func ExtractAndShowProgress(
	zip Zip,
	multiProgressBar interfaces.ProgressBarAdder,
) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		doneChan           = make(chan bool, 1)
		oopsBuilder        = oops.
					Code("ExtractAndShowProgress").
					With("zip.Name", zip.Name).
					With("zip.AbsPath", zip.AbsPath).
					With("multiProgressBar", multiProgressBar)

		extractedBytes int
		done           bool
		zipSize        int64
	)

	if zipSize, err = utils_zip.GetZipSize(zip.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", zip.AbsPath)
		return
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		zip.SetExtractionBar(multiProgressBar, zipSize)

		if err = utils_zip.Unzip(zip.AbsPath, filepath.Dir(zip.AbsPath), extractedBytesChan); err != nil {
			err = oopsBuilder.
				With("absPath", zip.AbsPath).
				With("extractedBytesChan", extractedBytesChan).
				Wrapf(err, "failed to unzip: %v", zip.AbsPath)
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
				zip.IncrementExtractionBar(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", zip.Name)
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
