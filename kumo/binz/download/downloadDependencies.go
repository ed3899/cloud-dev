package download

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/binz/download/progressBar"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type Binary struct {
	Dependency *draft.Dependency
	Err        error
}

type Binaries struct {
	Packer    *Binary
	Terraform *Binary
}

// Downloads dependencies and returns a struct containing the binaries
// If dps is empty, returns the already downloaded binaries
func DownloadDependencies(dps *draft.Dependencies) (*Binaries, error) {
	// Initiate download of dependencies
	log.Printf("Downloading %d dependencies...\n", len(*dps))

	// Create a channel to receive download results
	downloads := make(chan *progressBar.DownloadResult, len(*dps))
	// Create a channel to receive unzip results
	binariesChan := make(chan *Binary, 2)
	// Create a channel to receive errors
	errChan := make(chan error, 2)

	// Create a wait group to wait for all downloads to complete
	wg := sync.WaitGroup{}
	// Add 2 to the wait group for each dependency (1 for download, 1 for unzip)
	wg.Add(len(*dps) * 2)
	// Create a wait group to wait for all bars to complete
	bwg := sync.WaitGroup{}
	// Add 1 to the wait group for each bar
	progress := mpb.New(mpb.WithWaitGroup(&bwg), mpb.WithWidth(60), mpb.WithAutoRefresh())

	// Close channels when pipeline is complete
	go func() {
		wg.Wait()
		close(downloads)
		close(binariesChan)
		close(errChan)
	}()

	// Downloading...
	bwg.Add(1)
	go func(dps *draft.Dependencies, p *mpb.Progress) {
		defer bwg.Done()
		// Range over dependencies
		for _, dep := range *dps {
			// Download the dependency
			go func(dep *draft.Dependency, p *mpb.Progress) {
				defer wg.Done()
				progressBar.AttachDownloadBar(p, dep)
				DownloadDependency(dep, downloads)
			}(dep, progress)
		}

	}(dps, progress)

	// Unzipping...
	bwg.Add(1)
	go func(errChan chan<- error, progress *mpb.Progress) {
		defer bwg.Done()
		// Range over downloads channel
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

			// Unzip the download
			go func(dr *progressBar.DownloadResult, p *mpb.Progress) {
				defer wg.Done()
				progressBar.AttachZipBar(p, dr)
				Unzip(dr, binariesChan)
			}(dr, progress)
		}
	}(errChan, progress)

	// Wait for bars to complete and flush
	bwg.Wait()
	progress.Shutdown()

	// Check for errors
	if err, ok := <-errChan; ok {
		return nil, err
	}

	// Assemble binaries
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
		case "terraform":
			binaries.Terraform = binary
			continue
		}
	}

	fmt.Println("All dependencies downloaded!")

	return binaries, nil
}
