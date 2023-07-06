package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
)

func HandleDependencies() {
	// host := GetHostSpecs()
	// validHost := ValidateHostCompatibility(host)
	// packerDependency := DraftPackerDependency(validHost)
	// pulumiDependency := DraftPulumiDependency(validHost)
	// dependencies := []*Dependency{packerDependency, pulumiDependency}
	// barsWg := AppendBarToDependencies(dependencies)

	// resultChan := make(chan DownloadResult, len(dependencies))
	// dependenciesChan := make(chan *Dependency, len(dependencies))

	// go DownloadDependencies2(barsWg, dependencies, dependenciesChan, resultChan)

	// for _, v := range dependencies {
	// 	dependenciesChan <- v
	// }

	// DownloadDependencies(dependencies)
}

// func DownloadDependencies2(wg *sync.WaitGroup, deps []*Dependency, depsChan <-chan *Dependency, resultChan chan<- DownloadResult) {

// 	// TODO Try adding wg.add here
// 	for dep := range depsChan {
// 		go HandleDownload(dep, resultChan)
// 	}

// 	// go func() {
// 	// 	wg.Wait()

// 	// 	close(resultChan)
// 	// }()
// }

func DownloadDependencies(deps []*Dependency) {
	resultChan := make(chan DownloadResult, len(deps))
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))
	bars := GenerateBars(progress, deps)
	wg.Add(len(deps))

	for i, v := range deps {
		go func(dep *Dependency, bar *mpb.Bar) {
			defer wg.Done()
			if dep.Present {
				log.Printf("Dependency '%s' already exists. Skipped", filepath.Base(dep.ZipPath))
				return
			}
			HandleDownloads(dep, resultChan, bar)
		}(v, bars[i])
	}

	go func() {
		wg.Wait()

		close(resultChan)
	}()

	for result := range resultChan {
		if result.Err != nil {
			log.Printf("Error downloading %s: %v\n", result.Dependency.URL, result.Err)
		}
	}
}

// func HandleDownload(dep *Dependency, resultChan chan<- DownloadResult) {
// 	err := Download2(dep)

// 	result := DownloadResult{
// 		Dependency: dep,
// 		Err:        err,
// 	}

// 	resultChan <- result
// }

func HandleDownloads(dep *Dependency, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := Download(dep.URL, dep.ZipPath, bar)

	result := DownloadResult{
		Dependency: dep,
		Err:        err,
	}

	resultChan <- result
}

func Download2(dep *Dependency, downloads chan<- *DownloadResult, wg *sync.WaitGroup) {
	// Download
	defer wg.Done()
	url := dep.URL
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s'", url)
		log.Printf("error: %#v", err)
	}
	defer response.Body.Close()
	zipPath := dep.ZipPath
	destDir := filepath.Dir(zipPath)

	// Create the file along with all the necessary directories
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Printf("there was an error while creating %#v", destDir)
		log.Printf("error: %#v", err)
		return
	}

	// Create
	file, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Printf("there was an error while creating %#v", zipPath)
		log.Printf("error: %#v", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 4096)

	for {
		bytesDownloaded, err := response.Body.Read(buffer)

		if err != nil && err != io.EOF {
			log.Printf("err while downloading from '%s': %#v", url, err)
		}

		if bytesDownloaded == 0 {
			break
		}

		dep.Bar.IncrBy(bytesDownloaded)

		_, err = file.Write(buffer[:bytesDownloaded])

		if err != nil {
			log.Printf("there was an error while writing to %#v", zipPath)
			log.Printf("error: %#v", err)
		}

	}

	download := &DownloadResult{
		Dependency: dep,
		Err:        err,
		Fulfilled:  true,
	}

	downloads <- download
}

func Download(url string, path string, bar *mpb.Bar) error {
	// Download
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}
	defer response.Body.Close()
	destDir := filepath.Dir(path)

	// Create the file along with all the necessary directories
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Printf("there was an error while creating %#v", destDir)
		return err
	}

	// Create
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Printf("there was an error while creating %#v", path)
		return err
	}
	defer file.Close()

	buffer := make([]byte, 4096)

	for {
		bytesDownloaded, err := response.Body.Read(buffer)

		if err != nil && err != io.EOF {
			log.Printf("err while downloading from '%s': %#v", url, err)
		}

		if bytesDownloaded == 0 {
			break
		}

		bar.IncrBy(bytesDownloaded)

		_, err = file.Write(buffer[:bytesDownloaded])

		if err != nil {
			log.Printf("there was an error while writing to %#v", path)
			return err
		}

	}

	return nil
}
