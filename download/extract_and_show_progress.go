package download

import (
	"github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func ExtractAndShowProgress(
	progress *mpb.Progress,
	download *Download,
) error {
	oopsBuilder := oops.
		Code("ExtractAndShowProgress").
		With("download", download).
		With("progress", progress)

	extractedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	zipSize, err := zip.GetZipSize(download.Path.Zip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", download.Path.Zip)

		return err
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		download.Bar.Extracting = progress.AddBar(zipSize,
			mpb.BarQueueAfter(download.Bar.Downloading),
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

		err = zip.Unzip(download.Path.Zip, download.Path.Executable, extractedBytesChan)
		if err != nil {
			err := oopsBuilder.
				With("extractedBytesChan", extractedBytesChan).
				With("download.Path.Executable", download.Path.Executable).
				Wrapf(err, "failed to unzip: %v", download.Path.Zip)

			errChan <- err

			return
		}

		doneChan <- true
	}(zipSize)

OuterLoop:
	for {
		select {
		case extractedBytes := <-extractedBytesChan:
			if extractedBytes > 0 {
				download.Bar.Extracting.IncrBy(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", download.Name)

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
