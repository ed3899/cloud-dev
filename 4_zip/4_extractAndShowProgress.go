package zip

import (
	utils "github.com/ed3899/kumo/1_utils"
	"github.com/samber/oops"
)

func ExtractAndShowProgress(
	zip Zip,
	absPathToExtraction string,
	multiProgressBar MultiProgressBar,
) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		doneChan           = make(chan bool, 1)
		absPathToZip       = zip.AbsPath
		oopsBuilder        = oops.
					Code("extract_and_show_progress_failed").
					With("absPathToExtraction", absPathToExtraction).
					With("e.GetName()", zip.Name).
					With("e.GetPath()", zip.AbsPath).
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

		zip.SetExtractionBar(multiProgressBar, zipSize)

		if err := zip.ExtractTo(absPathToExtraction, extractedBytesChan); err != nil {
			err = oopsBuilder.
				With("extractedBytesChan", extractedBytesChan).
				Wrapf(err, "Error occurred while extracting %s", zip.Name)
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
