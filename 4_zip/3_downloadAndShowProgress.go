package zip

import (
	"github.com/samber/oops"
)

func DownloadAndShowProgress(
	z Zip,
	multiProgressBar MultiProgressBar,
) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		doneChan            = make(chan bool, 1)
		oopsBuilder         = oops.
					Code("download_and_show_progress_failed").
					With("d.Name", z.Name).
					With("d.AbsPath", z.AbsPath).
					With("multiProgressBar", multiProgressBar)

		downloadedBytes int
		done            bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		z.SetDownloadBar(multiProgressBar)

		if err := z.Download(downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "Error occurred while downloading %s", z.Name)
			errChan <- err
			return
		}

		doneChan <- true
	}()

OuterLoop:
	for {
		select {
		case downloadedBytes = <-downloadedBytesChan:
			if downloadedBytes > 0 {
				z.IncrementDownloadBar(downloadedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", z.Name)
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
