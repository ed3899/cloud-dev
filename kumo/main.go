package main

import (
	"fmt"
	"sync"

	"github.com/ed3899/kumo/utils"
)

func init() {
	host := utils.GetHostSpecs()
	validHost := utils.ValidateHostCompatibility(host)
	packerDependency := utils.DraftPackerDependency(validHost)
	pulumiDependency := utils.DraftPulumiDependency(validHost)
	dependencies := []*utils.Dependency{packerDependency, pulumiDependency}
	downloads := make(chan *utils.DownloadResult, len(dependencies))

	wg := sync.WaitGroup{}
	wg.Add(len(dependencies))
	utils.AppendBarToDependencies(&wg, dependencies)

	for _, dep := range dependencies {
		go utils.Download2(dep, downloads, &wg)
	}

	// Start a goroutine to wait for all downloads to complete
	go func() {
		wg.Wait()        // Wait until all downloads are complete
		close(downloads) // Close the 'downloads' channel to exit the loop
	}()

	for download := range downloads {
		fmt.Println(download.Fulfilled, download.Err)
	}

	// wg.Wait()   // Wait until all downloads are complete
	// Close the 'done' channel to indicate that we're done

	fmt.Println("All files downloaded!")
}

func main() {

}
