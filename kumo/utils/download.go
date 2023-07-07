package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(dep *Dependency, downloads chan<- *DownloadResult) error {
	url := dep.URL
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	zipPath := dep.ZipPath
	destDir := filepath.Dir(zipPath)

	// Create the destination along with all the necessary directories
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// Create file to write to
	file, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 4096)

	// Iterate over the response body and write to the file while updating the progress bar
	for {
		bytesDownloaded, err := response.Body.Read(buffer)

		if err != nil && err != io.EOF {
			return err
		}

		if bytesDownloaded == 0 {
			break
		}

		dep.DownloadBar.IncrBy(bytesDownloaded)

		_, err = file.Write(buffer[:bytesDownloaded])

		if err != nil {
			return err
		}

	}

	// Create the download result and send it to the channel
	download := &DownloadResult{
		Dependency: dep,
		Err:        nil,
	}

	downloads <- download

	return nil
}
