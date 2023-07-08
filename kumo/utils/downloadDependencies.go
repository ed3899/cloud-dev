package utils

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func DownloadDependencies(dps *Dependencies) (*Binaries, error) {
	// If there are no dependencies to be downloaded, return
	if len(*dps) == 0 {
		log.Println("All dependencies are installed")
		return nil, nil
	}

	// Create a channel to receive download results
	downloads := make(chan *DownloadResult, len(*dps))

	// Create a wait group to wait for all downloads to complete
	wg := sync.WaitGroup{}
	bwg := sync.WaitGroup{}

	// Add 2 to the wait group for each dependency (1 for download, 1 for unzip)
	wg.Add(len(*dps) * 2)
	progress := mpb.New(mpb.WithWaitGroup(&bwg), mpb.WithWidth(60), mpb.WithAutoRefresh())

	// Start a download for each dependency
	bwg.Add(1)
	go func(dps *Dependencies, p *mpb.Progress) {
		defer bwg.Done()

		for _, dep := range *dps {
			go func(dep *Dependency, p *mpb.Progress) {
				defer wg.Done()
				AttachDownloadBar(p, dep)
				Download(dep, downloads)
			}(dep, progress)
		}

	}(dps, progress)

	// Create a channel to receive unzip results
	binariesChan := make(chan *Binary, 2)
	errChan := make(chan error, 2)

	go func() {
		wg.Wait()
		close(downloads)
		close(binariesChan)
		close(errChan)
	}()

	bwg.Add(1)
	go func(errChan chan<- error, progress *mpb.Progress) {
		defer bwg.Done()
		// Start a goroutine to unzip each dependency
		for dr := range downloads {
			if dr.Err != nil {
				// Remove the download if there was an error
				msg := fmt.Sprintf("Error occurred while downloading %s", dr.Dependency.Name)
				err := errors.Wrap(dr.Err, msg)
				log.Println(err)

				log.Printf("Removing failed download %s...\n", dr.Dependency.ZipPath)
				err = os.RemoveAll(dr.Dependency.ZipPath)
				// If there was an error removing the failed download, return
				if err != nil {
					msg := fmt.Sprintf("Error occurred while removing failed download %s", dr.Dependency.ZipPath)
					err := errors.Wrap(err, msg)
					errChan <- err
					return
				}

				continue
			}

			go func(dr *DownloadResult, p *mpb.Progress) {
				defer wg.Done()
				AttachZipBar(p, dr)
				Unzip(dr, binariesChan)
			}(dr, progress)
		}
	}(errChan, progress)

	bwg.Wait()

	if err, ok := <-errChan; ok {
		return nil, err
	}

	// Start a goroutine to wait for all binaries to be created
	binaries := &Binaries{}

	for binary := range binariesChan {
		if binary.Err != nil {
			msg := fmt.Sprintf("Error occurred while unzipping %s", binary.Dependency.Name)
			err := errors.Wrap(binary.Err, msg)
			log.Println(err)

			log.Printf("Removing failed unzip %s...\n", binary.Dependency.ExtractionPath)
			err = os.RemoveAll(binary.Dependency.ExtractionPath)
			// If there was an error removing the failed unzip, return
			if err != nil {
				msg := fmt.Sprintf("Error occurred while removing failed unzip %s", binary.Dependency.ExtractionPath)
				err := errors.Wrap(err, msg)
				log.Println(err)
				return nil, err
			}

			continue
		}

		switch binary.Dependency.Name {
		case "packer":
			binaries.Packer = binary
			continue
		case "pulumi":
			binaries.Pulumi = binary
			continue
		}
	}

	fmt.Println("All dependencies downloaded!")

	return binaries, nil
}
