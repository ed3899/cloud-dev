package binaries

import (
	"os"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type ExtractableByWorkflow interface {
	Extractable
	Removable
	Retrivable
}

func ExtractAndShowProgress[E ExtractableByWorkflow](e E, multiProgressBar *mpb.Progress) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		done               = make(chan bool, 1)
		zipSize            int64
		absPathToZip       = e.GetPath()
	)

	if zipSize, err = utils.GetZipSize(absPathToZip); err != nil {
		err = errors.Wrapf(err, "failed to get zip size for: %v", absPathToZip)
		return
	}

	go func(zipSize int64) {
		defer close(extractedBytesChan)
		defer close(errChan)
		defer close(done)

		e.SetExtractionBar(multiProgressBar, zipSize)
		if err = e.Extract(absPathToZip, extractedBytesChan); err != nil {
			errChan <- err
			return
		}
		done <- true
	}(zipSize)

OuterLoop:
	for {
		select {
		case extractedBytes := <-extractedBytesChan:
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
				break OuterLoop
			}

			break OuterLoop

		case d := <-done:
			if d {
				break OuterLoop
			}
			continue OuterLoop
		}
	}

	return
}
