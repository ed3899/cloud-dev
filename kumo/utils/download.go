package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Download(dep *Dependency, downloads chan<- *DownloadResult) {
	url := dep.URL
	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrap(err, "failed to get dependency")
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
		err = errors.Wrap(err, "failed to create destination directory")
		downloads <- &DownloadResult{
			Dependency: dep,
			Err:        err,
		}
		return
	}

	// Create file to write to
	file, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		err = errors.Wrap(err, "failed to create file")
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
			err = errors.Wrap(err, "failed to read response body")
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
			err = errors.Wrap(err, "failed to write to file")
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
