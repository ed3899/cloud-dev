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
	"github.com/vbauerster/mpb/v8/decor"
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

	// Create the destination along with all the necessary directories
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

	// Create file to write to
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

	// Iterate over the response body and write to the file while updating the progress bar
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

	// Create the download result and send it to the channel
	download := &DownloadResult{
		Dependency: dep,
		Err:        nil,
	}

	downloads <- download
}

func attachProgressBar(wg *sync.WaitGroup, d *Dependency) *Dependency {
	progress := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))

	downloadBar := progress.AddBar(int64(d.ContentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(d.Name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"done",
			),
		),
	)

	d.DownloadBar = downloadBar

	return d
}
