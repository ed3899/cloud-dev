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
	wg.Add(len(dependencies) * 2)
	progress := utils.AppendDownloadBar(&wg, dependencies)

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

	for dr := range downloads {
		if dr.Err != nil {
			fmt.Printf("Error occurred while downloading %s: %v\n", dr.Dependency.Name, dr.Err)
			continue
		}

		utils.AppendZipBar(progress, dr)
		go func(dr *utils.DownloadResult) {
			defer wg.Done()
			err := utils.UnzipSource(dr)
			if err != nil {
				fmt.Printf("Error occurred while unzipping %s: %v\n", dr.Dependency.Name, err)
				return
			}
		}(dr)
	}

	fmt.Println("All files downloaded!")
}

func main() {

}
