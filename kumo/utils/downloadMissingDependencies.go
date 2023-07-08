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
	wg.Add(len(*dps) * 3)
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))
	AttachDownloadBar(progress, *dps)

	// Start a download for each dependency
	for _, dep := range *dps {
		go func(dep *Dependency) {
			defer wg.Done()
			err := Download(dep, downloads)
			if err != nil {
				downloads <- &DownloadResult{
					Dependency: dep,
					Err:        err,
				}
				return
			}
		}(dep)
	}

	// Start a goroutine to wait for all downloads to complete
	go func() {
		wg.Wait()
		close(downloads)
	}()

	// Create a channel to receive unzip results
	binaries := make(chan *Binary, 2)

	// Start a goroutine to unzip each dependency
	for dr := range downloads {
		if dr.Err != nil {
			fmt.Printf("Error occurred while downloading %s: %v\n", dr.Dependency.Name, dr.Err)
			continue
		}

		AttachZipBar(progress, dr)

		go func(dr *DownloadResult) {
			defer wg.Done()
			err := Unzip(dr, binaries)
			if err != nil {
				fmt.Printf("Error occurred while unzipping %s: %v\n", dr.Dependency.Name, err)
				return
			}
		}(dr)
	}

	// Start a goroutine to wait for all unzips to complete and flush the progress bar
	go func() {
		progress.Wait()
	}()


	// Start a goroutine to wait for all binaries to be created
	go func() {
		
	}()


	wg.Wait()

	fmt.Println("All dependencies downloaded!")

	return &Binaries{}, nil
}
