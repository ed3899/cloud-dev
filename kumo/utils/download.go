package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Download(url, destPath string, bytesDownloadedChan chan<- int) (err error) {
	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to download from: %s", url)
		return
	}
	defer response.Body.Close()

	destDir := filepath.Dir(destPath)

	// Create the destination along with all the necessary directories
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		err = errors.Wrapf(err, "failed to create destination directory for: %s", destPath)
		return
	}

	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		err = errors.Wrapf(err, "failed to create file for: %s", destPath)
		return
	}
	defer file.Close()

	buffer := make([]byte, 4096)

	// Iterate over the response body and sned downloaded bytes to channel
	for {
		bytesDownloaded, err := response.Body.Read(buffer)

		if err != nil && err != io.EOF {
			err = errors.Wrap(err, "failed to read response body")
			return err
		}

		if bytesDownloaded == 0 {
			break
		}

		bytesDownloadedChan <- bytesDownloaded

		_, err = file.Write(buffer[:bytesDownloaded])
		if err != nil {
			err = errors.Wrap(err, "failed to write to file")
			return err
		}
	}

	return
}
