package download

import (
	common_zip_interfaces "github.com/ed3899/kumo/common/zip/interfaces"
	"github.com/samber/oops"
)

type DownloadableAndProgressive interface {
	common_zip_interfaces.Downloadable
	common_zip_interfaces.Removable
	common_zip_interfaces.Retrivable
}

func DownloadAndShowProgress(
	d DownloadableAndProgressive,
	multiProgressBar common_zip_interfaces.MultiProgressBar,
) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		doneChan            = make(chan bool, 1)
		oopsBuilder         = oops.
					Code("download_and_show_progress_failed").
					With("d.GetName()", d.GetName()).
					With("d.GetPath()", d.GetPath()).
					With("multiProgressBar", multiProgressBar)

		downloadedBytes int
		done            bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		d.SetDownloadBar(multiProgressBar)

		if err := d.Download(downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "Error occurred while downloading %s", d.GetName())
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
				d.IncrementDownloadBar(downloadedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", d.GetName())
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
