package download

import (
	"github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func ExtractAndShowProgress(
	progress *mpb.Progress,
	download *Download,
) {
	oopsBuilder := oops.
		Code("ExtractAndShowProgress").
		With("download", download).
		With("progress", progress)

	extractedBytesChan := make(chan int, 1024)
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	zipSize, err := url.GetContentLength()
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", download.Path.Zip)
		return
	}

}
