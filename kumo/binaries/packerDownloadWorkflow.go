package binaries

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func PackerDownloadWorkflow(packer *Packer) (err error) {
	var (
		progressWaitGroup *sync.WaitGroup
		multiProgressBar  = mpb.New(mpb.WithWaitGroup(progressWaitGroup), mpb.WithWidth(60), mpb.WithAutoRefresh())
	)

	if err = DownloadAndShowProgress(packer.Zip, multiProgressBar); err != nil {
		err = errors.Wrapf(err, "failed to download %s", packer.Zip.Name)
		return
	}

	if err = ExtractAndShowProgress(packer.Zip, multiProgressBar); err != nil {
		err = errors.Wrapf(err, "failed to extract %s", packer.Zip.Name)
		return
	}

	multiProgressBar.Shutdown()

	return
}
