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

func DownloadDependencies(ze []*Dependency) {
	resultChan := make(chan DownloadResult, len(ze))
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))
	bars := GenerateBars(progress, ze)
	wg.Add(len(ze))

	for i, v := range ze {
		go func(ref *Dependency, bar *mpb.Bar) {
			defer wg.Done()
			if ref.Present {
				log.Printf("File '%s' already exists", ref.ExtractionPath)
				return
			}
			DownloadZip(ref, resultChan, bar)
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

func DownloadZip(ref *Dependency, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := Download(ref.URL, ref.ExtractionPath, bar)

	result := DownloadResult{
		Dependency: ref,
		Err:        err,
	}

	resultChan <- result
}

func Download(url string, binPath string, bar *mpb.Bar) error {
	// Download
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}
	defer response.Body.Close()

	// Create the file along with all the necessary directories
	err = os.MkdirAll(filepath.Dir(binPath), 0755)
	if err != nil {
		log.Printf("there was an error while creating %#v", binPath)
		return err
	}

	// Create
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Printf("there was an error while creating %#v", binPath)
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
			log.Printf("there was an error while writing to %#v", binPath)
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

