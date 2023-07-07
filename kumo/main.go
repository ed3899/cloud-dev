package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ed3899/kumo/utils"
	"github.com/vbauerster/mpb/v8"
)

func init() {
	host := utils.GetHostSpecs()
	validHost := utils.ValidateHostCompatibility(host)
	packerDependency := utils.DraftPackerDependency(validHost)
	pulumiDependency := utils.DraftPulumiDependency(validHost)
	dependencies := []*utils.Dependency{packerDependency, pulumiDependency}
	downloads := make(chan *utils.DownloadResult, len(dependencies))

	wg := sync.WaitGroup{}
	wg.Add(len(dependencies) * 2)
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))
	utils.AttachDownloadBar(progress, dependencies)

	// Start a download for each dependency
	for _, dep := range dependencies {
		go func(dep *utils.Dependency) {
			defer wg.Done()
			err := utils.Download(dep, downloads)
			if err != nil {
				downloads <- &utils.DownloadResult{
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

	// Start a goroutine to unzip each dependency
	for dr := range downloads {
		if dr.Err != nil {
			fmt.Printf("Error occurred while downloading %s: %v\n", dr.Dependency.Name, dr.Err)
			continue
		}

		utils.AttachZipBar(progress, dr)

		go func(dr *utils.DownloadResult) {
			defer wg.Done()
			err := utils.Unzip(dr)
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

	wg.Wait()

	fmt.Println("All files downloaded!")
}

func main() {

}
