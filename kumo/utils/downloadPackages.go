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

func HandleDownloads(dep *Dependency, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := Download(dep.URL, dep.ZipPath, bar)

	result := DownloadResult{
		Dependency: dep,
		Err:        err,
	}

	resultChan <- result
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

	// Fill
	// _, err = io.Copy(file, response.Body)
	// if err != nil {
	// 	log.Printf("there was an error while copying contents to %#v", binPath)
	// 	return err
	// }

	return nil
}

