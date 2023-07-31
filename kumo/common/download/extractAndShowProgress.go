package download

import (
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type ExtractableAndProgressive interface {
	Extractable
	Removable
	Retrivable
}

func ExtractAndShowProgress[E ExtractableAndProgressive](e E, absPathToExtraction string, multiProgressBar MultiProgressBarI) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		doneChan           = make(chan bool, 1)
		absPathToZip       = e.GetPath()
		oopsBuilder        = oops.
					Code("extract_and_show_progress_failed").
					With("absPathToExtraction", absPathToExtraction).
					With("e.GetName()", e.GetName()).
					With("e.GetPath()", e.GetPath()).
					With("multiProgressBar", multiProgressBar)

		extractedBytes int
		done           bool
		zipSize        int64
	)

	if zipSize, err = utils.GetZipSize(absPathToZip); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", absPathToZip)
		return
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		e.SetExtractionBar(multiProgressBar, zipSize)

		if err := e.ExtractTo(absPathToExtraction, extractedBytesChan); err != nil {
			err = oopsBuilder.
				With("extractedBytesChan", extractedBytesChan).
				Wrapf(err, "Error occurred while extracting %s", e.GetName())
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
				e.IncrementExtractionBar(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", e.GetName())
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
