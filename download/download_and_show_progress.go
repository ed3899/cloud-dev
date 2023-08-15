package download

import (
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

type IDownloadUrlAndPath interface {
	IUrlGetter
	IPathGetter
}

func DownloadAndShowProgressWith(
	bar *mpb.Bar,
	utilsUrlDownload func(download IDownloadUrlAndPath, downloadedBytesChan chan<- int) error,
) DownloadAndShowProgress {
	oopsBuilder := oops.
		Code("DownloadAndShowProgressWith")

	// var wg sync.WaitGroup
	// mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh(), mpb.WithWidth(64))
	downloadedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	downloadAndShowProgress := func(download IDownload) (Download, error) {
		downloadClone := download.Clone()

		go func() {
			defer close(downloadedBytesChan)
			defer close(errChan)
			defer close(doneChan)

			downloadClone.Bar().SetDownloading(bar)

			err := utilsUrlDownload(downloadClone, downloadedBytesChan)
			if err != nil {
				err = oopsBuilder.
					With("path", downloadClone.Path()).
					With("downloadedBytesChan", downloadedBytesChan).
					Wrapf(err, "failed to download: %v", downloadClone.Url())
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
					downloadClone.Bar().Downloading().IncrBy(downloadedBytes)
				}

			case err := <-errChan:
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "Error occurred while downloading %s", downloadClone.Name())
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

	return downloadAndShowProgress
}

type DownloadAndShowProgress func(IDownload) (Download, error)
