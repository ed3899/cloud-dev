package zip

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func DownloadAndShowProgress(
	download *Download,
	multiProgressBar interfaces.MpbV8MultiprogressBar,
	urlDownload url.DownloadF,
) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		doneChan            = make(chan bool, 1)
		oopsBuilder         = oops.
					Code("download_and_show_progress_failed").
					With("download.Name", download.Name).
					With("download.AbsPath", download.AbsPath).
					With("multiProgressBar", multiProgressBar)

		downloadedBytes int
		done            bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		download.SetDownloadBar(multiProgressBar)

		if err = urlDownload(download.Url, download.AbsPath, downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("url", download.Url).
				With("absPath", download.AbsPath).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", download.Url)
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
				download.IncrementDownloadBar(downloadedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", download.Name)
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
