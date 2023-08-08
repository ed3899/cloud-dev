package zip

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func DownloadAndShowProgress(
	zip Zip,
	multiProgressBar interfaces.ProgressBarAdder,
) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		doneChan            = make(chan bool, 1)
		oopsBuilder         = oops.
					Code("download_and_show_progress_failed").
					With("zip.Name", zip.Name).
					With("zip.AbsPath", zip.AbsPath).
					With("multiProgressBar", multiProgressBar)

		downloadedBytes int
		done            bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		zip.SetDownloadBar(multiProgressBar)

		if err = url.Download(zip.Url, zip.AbsPath, downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("url", zip.Url).
				With("absPath", zip.AbsPath).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", zip.Url)
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
				zip.IncrementDownloadBar(downloadedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", zip.Name)
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
