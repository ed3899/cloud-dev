package binaries

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func ExtractWorkflow[Z ZipI](z Z, multiProgressBar *mpb.Progress) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		done               = make(chan bool, 1)
	)

	go func() {
		defer close(extractedBytesChan)
		defer close(errChan)
		defer close(done)

		z.SetExtractionBar(multiProgressBar)
		if err = z.Extract(z.GetPath(), extractedBytesChan); err != nil {
			errChan <- err
			return
		}
		done <- true
	}()

OuterLoop:
	for {
		select {
		case extractedBytes := <-extractedBytesChan:
			if extractedBytes <= 0 {
				continue OuterLoop
			}

			if err = z.IncrementExtractionBar(extractedBytes); err != nil {
				log.Print(err)
				continue OuterLoop
			}
		case err = <-errChan:
			if err == nil {
				continue OuterLoop
			}

			if err := os.RemoveAll(z.GetPath()); err != nil {
				err = errors.Wrapf(err, "Error occurred while removing %s", z.GetName())
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
