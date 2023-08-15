package download

import (
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndShowProgress(
	bar *mpb.Bar,
	download *Download,
) (*Download, error) {
	oopsBuilder := oops.
		Code("DownloadAndShowProgress")

	// var wg sync.WaitGroup
	// mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh(), mpb.WithWidth(64))
	downloadedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	downloadClone := download.Clone()

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		downloadClone.Bar.Downloading = bar

		err := url.Download(downloadClone.Url, downloadClone.Path.Zip, downloadedBytesChan)
		if err != nil {
			err = oopsBuilder.
				With("path", downloadClone.Path).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", downloadClone.Url)
			errChan <- err
			return
		}

		doneChan <- true
	}()

OuterLoop:
	for {
		select {
		case downloadedBytes := <-downloadedBytesChan:
			if downloadedBytes > 0 {
				downloadClone.Bar.Downloading.IncrBy(downloadedBytes)
			}

		case err := <-errChan:
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", downloadClone.Name)
				return downloadClone, err
			}

		case done := <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return downloadClone, nil
}
