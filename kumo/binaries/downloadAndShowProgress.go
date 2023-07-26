package binaries

import (
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type DownloadableByWorkflow interface {
	Downloadable
	Removable
	Retrivable
}

func DownloadAndShowProgress[D DownloadableByWorkflow](d D, multiProgressBar *mpb.Progress) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		downloadedBytes     int

		errChan = make(chan error, 1)

		doneChan = make(chan bool, 1)
		done     bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		d.SetDownloadBar(multiProgressBar)
		if err = d.Download(downloadedBytesChan); err != nil {
			errChan <- err
			return
		}
		doneChan <- true
	}()

OuterLoop:
	for {
		select {
		case downloadedBytes = <-downloadedBytesChan:
			if downloadedBytes <= 0 {
				continue OuterLoop
			}

			d.IncrementDownloadBar(downloadedBytes)

		case err = <-errChan:
			if err == nil {
				continue OuterLoop
			}

			if err = d.Remove(); err != nil {
				err = errors.Wrapf(err, "Error occurred while removing %s", d.GetName())
			}

			err = errors.Wrapf(err, "Error occurred while downloading %s", d.GetName())
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
