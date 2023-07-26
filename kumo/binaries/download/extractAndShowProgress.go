package download

import (
	"os"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type ExtractableAndProgressive interface {
	binaries.Extractable
	binaries.Removable
	binaries.Retrivable
}

func ExtractAndShowProgress[E ExtractableAndProgressive](e E, absPathToExtraction string, multiProgressBar *mpb.Progress) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		extractedBytes     int

		errChan = make(chan error, 1)

		doneChan     = make(chan bool, 1)
		done         bool
		zipSize      int64
		absPathToZip = e.GetPath()
	)

	if zipSize, err = utils.GetZipSize(absPathToZip); err != nil {
		err = errors.Wrapf(err, "failed to get zip size for: %v", absPathToZip)
		return
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		e.SetExtractionBar(multiProgressBar, zipSize)
		if err = e.ExtractTo(absPathToExtraction, extractedBytesChan); err != nil {
			errChan <- err
			return
		}
		doneChan <- true
	}(zipSize)

OuterLoop:
	for {
		select {
		case extractedBytes = <-extractedBytesChan:
			if extractedBytes <= 0 {
				continue OuterLoop
			}

			e.IncrementExtractionBar(extractedBytes)

		case err = <-errChan:
			if err == nil {
				continue OuterLoop
			}

			if err = os.RemoveAll(e.GetPath()); err != nil {
				err = errors.Wrapf(err, "Error occurred while removing %s", e.GetName())
			}

			err = errors.Wrapf(err, "Error occurred while extracting %s", e.GetName())
			break OuterLoop

		case done = <-doneChan:
			if done {
				break OuterLoop
			}
			continue OuterLoop
		}
	}

	return
}
