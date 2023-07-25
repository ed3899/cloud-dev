package binaries

import (
	"log"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)


type DownloadableByWorkflow interface {
	Downloadable
	Removable
	Retrivable
}

func DownloadWorkflow[D DownloadableByWorkflow](d D, multiProgressBar *mpb.Progress) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		done                = make(chan bool, 1)
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(done)

		d.SetDownloadBar(multiProgressBar)
		if err = d.Download(downloadedBytesChan); err != nil {
			errChan <- err
			return
		}
		done <- true
	}()

OuterLoop:
	for {
		select {
		case downloadedBytes := <-downloadedBytesChan:
			if downloadedBytes <= 0 {
				continue OuterLoop
			}

			if err = d.IncrementDownloadBar(downloadedBytes); err != nil {
				log.Print(err)
				continue OuterLoop
			}
		case err = <-errChan:
			if err == nil {
				continue OuterLoop
			}

			if err = d.Remove(); err != nil {
				err = errors.Wrapf(err, "Error occurred while removing %s", d.GetName())
			}

			err = errors.Wrapf(err, "Error occurred while downloading %s", d.GetName())
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
