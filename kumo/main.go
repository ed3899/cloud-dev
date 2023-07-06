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
	utils.AppendBarToDependencies(&wg, dependencies)
	done := make(chan bool)

	for _, dep := range dependencies {
		wg.Add(1)

		go utils.Download2(dep, downloads, &wg)
	}

	wg.Wait()   // Wait until all downloads are complete
	close(done) // Close the 'done' channel to indicate that we're done

	fmt.Println("All files downloaded!")
}

func main() {

}
