package binaries

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func GetPackerWorkflow() (err error) {
	var (
		packer *Packer
		progressWaitGroup *sync.WaitGroup
		multiProgressBar = mpb.New(mpb.WithWaitGroup(progressWaitGroup), mpb.WithWidth(60), mpb.WithAutoRefresh())
	)

	if packer, err = NewPacker(); err != nil {
		err = errors.Wrap(err, "failed to create packer instance")
		return
	}

	if err = DownloadWorkflow(packer.Zip, multiProgressBar); err != nil {
		err = errors.Wrap(err, "failed to download packer")
		return
	}

	if err = ExtractWorkflow(packer.Zip, multiProgressBar); err != nil {
		err = errors.Wrap(err, "failed to extract packer")
		return
	}

	multiProgressBar.Shutdown()

	return
}