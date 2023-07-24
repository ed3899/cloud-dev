package binaries

import (
	"log"
	"sync"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndExtractWorkflow[Z ZipI](d Z) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		extractedBytesChan	= make(chan int, 1024)
		errChan             = make(chan error, 1)
		progressGroup       = &sync.WaitGroup{}
		progressBar         = mpb.New(mpb.WithWaitGroup(progressGroup), mpb.WithWidth(60), mpb.WithAutoRefresh())
	)

	progressGroup.Add(1)
	go func() {
		defer progressGroup.Done()
		defer close(downloadedBytesChan)
		defer close(errChan)

		d.SetDownloadBar(progressBar)
		if err = d.Download(downloadedBytesChan); err != nil {
			errChan <- err
			return
		}
	}()

	for {
		select {
		case downloadedBytes := <-downloadedBytesChan:
			if downloadedBytes <= 0 {
				break
			}

			if err = d.IncrementDownloadBar(downloadedBytes); err != nil {
				log.Print(err)
				continue
			}
		case err = <-errChan:
			if err = d.Remove(); err != nil {
				err = errors.Wrapf(err, "Error occurred while removing %s", d.GetName())
				return
			}
			continue
		default:
			continue
		}
	}

	// progressGroup.Add(1)
	// go func() {
	// 	defer progressGroup.Done()
	// 	defer close(extractedBytesChan)
	// 	defer close(errChan)


	// }()

}
