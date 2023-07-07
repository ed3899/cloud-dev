package main

import (
	"fmt"
	"sync"
	// "time"

	"github.com/ed3899/kumo/utils"
	// "github.com/vbauerster/mpb/v8"
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
	progress := utils.AppendDownloadBar(&wg, dependencies)

	// Start a goroutine for each dependency
	for _, dep := range dependencies {
		go utils.Download(dep, downloads, &wg)
	}

	// Start a goroutine to wait for all downloads to complete
	go func() {
		wg.Wait()        // Wait until all downloads are complete
		close(downloads) // Close the 'downloads' channel to exit the loop
	}()

	for d := range downloads {
		// progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))

		if d.Err != nil {
			fmt.Printf("Error occurred while downloading %s: %v\n", d.Dependency.Name, d.Err)
			continue
		}

		utils.AppendZipBar(progress, d)
		go utils.UnzipSource(d, &wg)
	}

	fmt.Println("All files downloaded!")
}


func main() {

}
