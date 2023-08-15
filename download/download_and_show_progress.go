package download

import (
	"sync"

	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndShowProgressWith(
	bar *mpb.Bar,
	download IDownload,
	utilsUrlDownload func(url string, absPath string, downloadedBytesChan chan<- int) error,
) {
	oopsBuilder := oops.
		Code("DownloadAndShowProgressWith").
		With("download", download)

	var wg sync.WaitGroup
	mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh(), mpb.WithWidth(64))

	downloadClone := download.Clone()
	downloadedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		downloadClone.Bar().SetDownloading(bar)



		if err = urlDownload(dae.Download.Url, dae.Download.AbsPath, downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("url", dae.Download.Url).
				With("absPath", dae.Download.AbsPath).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", dae.Download.Url)
			errChan <- err
			return
		}

		doneChan <- true
	}()

}

type DownloadAndShowProgress func()
