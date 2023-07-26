package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// TODO add context
func Download(url, destPath string, bytesDownloadedChan chan<- int) (err error) {
	var (
		destDir     = filepath.Dir(destPath)
		bytesBuffer = make([]byte, 4096)

		response        *http.Response
		file            *os.File
		bytesDownloaded int
	)

	// Initiate download and defer closing the response body
	if response, err = http.Get(url); err != nil {
		return errors.Wrapf(err, "failed to download from: %s", url)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			err = errors.Wrap(err, "failed to close response body")
		}
	}()

	// Create the destination dir
	if err = os.MkdirAll(destDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create destination directory for: %s", destPath)
	}

	// Create the file to write to
	if file, err = os.Create(destPath); err != nil {
		return errors.Wrapf(err, "failed to create file for: %s", destPath)
	}
	defer file.Close()

	// Iterate over the response body
	for {
		// Read the response body into the bytes buffer
		if bytesDownloaded, err = response.Body.Read(bytesBuffer); err != nil && err != io.EOF {
			return errors.Wrap(err, "failed to read response body")
		}

		if bytesDownloaded == 0 {
			break
		}

		// Send the number of bytes downloaded to the channel
		bytesDownloadedChan <- bytesDownloaded

		// Write the bytes to the file
		if _, err = file.Write(bytesBuffer[:bytesDownloaded]); err != nil {
			return errors.Wrap(err, "failed to write to file")
		}
	}

	return
}
