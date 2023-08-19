package download

import (
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func (d *Download) DownloadAndShowProgress() error {
	oopsBuilder := oops.
		Code("DownloadAndShowProgress").
		In("download").
		Tags("Download")

	downloadedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		d.Bar.Downloading = d.Progress.AddBar(
			d.ContentLength,
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"unzipped",
				),
			),
		)

		err := url.Download(d.Url, d.Path.Zip, downloadedBytesChan)
		if err != nil {
			err = oopsBuilder.
				With("path", d.Path).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", d.Url)
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
				d.Bar.Downloading.IncrBy(downloadedBytes)
			}

		case err := <-errChan:
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", d.Name)

				return err
			}

		case done := <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return nil
}
