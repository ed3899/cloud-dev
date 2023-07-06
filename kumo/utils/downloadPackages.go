package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func Download(dep *Dependency, downloads chan<- *DownloadResult, wg *sync.WaitGroup) {
	// Download
	defer wg.Done()
	url := dep.URL
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s'", url)
		log.Printf("error: %#v", err)
		downloads <- &DownloadResult{
			Dependency: dep,
			Err:        err,
		}
		return
	}
	defer response.Body.Close()
	zipPath := dep.ZipPath
	destDir := filepath.Dir(zipPath)

	// Create the file along with all the necessary directories
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Printf("there was an error while creating %#v", destDir)
		log.Printf("error: %#v", err)
		downloads <- &DownloadResult{
			Dependency: dep,
			Err:        err,
		}
		return
	}

	// Create
	file, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Printf("there was an error while creating %#v", zipPath)
		log.Printf("error: %#v", err)
		downloads <- &DownloadResult{
			Dependency: dep,
			Err:        err,
		}
		return
	}
	defer file.Close()

	buffer := make([]byte, 4096)

	for {
		bytesDownloaded, err := response.Body.Read(buffer)

		if err != nil && err != io.EOF {
			log.Printf("err while downloading from '%s': %#v", url, err)
			downloads <- &DownloadResult{
				Dependency: dep,
				Err:        err,
			}
			return
		}

		if bytesDownloaded == 0 {
			break
		}

		dep.DownloadBar.IncrBy(bytesDownloaded)

		_, err = file.Write(buffer[:bytesDownloaded])

		if err != nil {
			log.Printf("there was an error while writing to %#v", zipPath)
			log.Printf("error: %#v", err)
			downloads <- &DownloadResult{
				Dependency: dep,
				Err:        err,
			}
			return
		}

	}

	download := &DownloadResult{
		Dependency: dep,
		Err:        nil,
	}

	downloads <- download
}
