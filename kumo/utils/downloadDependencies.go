package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

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

	// Add 3 to the wait group for each dependency (1 for download, 1 for unzip, 1 for attaching the binary)
	wg.Add(len(*dps) * 2)
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))
	AttachDownloadBar(progress, *dps)

	// Start a download for each dependency
	for _, dep := range *dps {
		go func(dep *Dependency) {
			defer wg.Done()
			Download(dep, downloads)
		}(dep)
	}

	// Create a channel to receive unzip results
	binariesChan := make(chan *Binary, 2)

	go func() {
		wg.Wait()
		close(downloads)
	}()

	// Start a goroutine to unzip each dependency
	for dr := range downloads {
		if dr.Err != nil {
			// Remove the download if there was an error
			log.Printf("Error occurred while downloading %s: %v\n", dr.Dependency.Name, dr.Err)
			continue
		}

		AttachZipBar(progress, dr)
		go func(dr *DownloadResult) {
			defer wg.Done()
			Unzip(dr, binariesChan)
		}(dr)
	}

	progress.Wait()
	close(binariesChan)

	// Start a goroutine to wait for all binaries to be created
	binaries := &Binaries{}

	for binary := range binariesChan {
		if binary.Err != nil {
			log.Printf("Error occurred while creating binary %s: %v\n", binary.Dependency.Name, binary.Err)
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
