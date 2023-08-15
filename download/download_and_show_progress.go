package download

import (
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func DownloadAndShowProgress(
	progress *mpb.Progress,
	download *Download,
) (*Download, error) {
	oopsBuilder := oops.
		Code("DownloadAndShowProgress")

	// var wg sync.WaitGroup
	// mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh(), mpb.WithWidth(64))
	downloadedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		download.Bar.Downloading = progress.AddBar(
			download.ContentLength,
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(download.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"unzipped",
				),
			),
		)

		err := url.Download(download.Url, download.Path.Zip, downloadedBytesChan)
		if err != nil {
			err = oopsBuilder.
				With("path", download.Path).
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
		case downloadedBytes := <-downloadedBytesChan:
			if downloadedBytes > 0 {
				download.Bar.Downloading.IncrBy(downloadedBytes)
			}

		case err := <-errChan:
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", download.Name)
				return download, err
			}

		case done := <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return download, nil
}
