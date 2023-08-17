package download

import (
	"path/filepath"

	"github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func (d *Download) ExtractAndShowProgress() error {
	oopsBuilder := oops.
		Code("ExtractAndShowProgress")

	extractedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	zipSize, err := zip.GetZipSize(d.Path.Zip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", d.Path.Zip)

		return err
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		d.Bar.Extracting = d.Progress.AddBar(zipSize,
			mpb.BarQueueAfter(d.Bar.Downloading),
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

		err = zip.Unzip(d.Path.Zip, filepath.Dir(d.Path.Executable), extractedBytesChan)
		if err != nil {
			err := oopsBuilder.
				With("extractedBytesChan", extractedBytesChan).
				With("download.Path.Executable", d.Path.Executable).
				Wrapf(err, "failed to unzip: %v", d.Path.Zip)

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
				d.Bar.Extracting.IncrBy(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", d.Name)

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
