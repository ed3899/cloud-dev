package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

// TODO add context
func Download(url, destPath string, bytesDownloadedChan chan<- int) (err error) {
	var (
		destDir     = filepath.Dir(destPath)
		bytesBuffer = make([]byte, 4096)
		oopsBuilder = oops.Code("download_failed").
				With("url", url).
				With("destPath", destPath).
				With("bytesDownloadedChan", bytesDownloadedChan).
				With("destDir", destDir).
				With("bytesBuffer", bytesBuffer)

		response        *http.Response
		downloadFile    *os.File
		bytesDownloaded int
		bytesWritten    int
	)

	// Initiate download and defer closing the response body
	if response, err = http.Get(url); err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to initiate download from: %s", url)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("responseStatus", response.Status).
					Wrapf(err, "failed to close response body"),
			)
		}
	}()

	// Create the destination dir
	if err = os.MkdirAll(destDir, 0755); err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to create destination directory for: %s", destPath)
	}

	// Create the file to write to
	if downloadFile, err = os.Create(destPath); err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to create file for: %s", destPath)
	}
	defer func() {
		if err = downloadFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("downloadFile", downloadFile).
					Wrapf(err, "failed to close downloadFile: %s", downloadFile.Name()),
			)
		}
	}()

	// Iterate over the response body
	for {
		// Read the response body into the bytes buffer
		if bytesDownloaded, err = response.Body.Read(bytesBuffer); err != nil && err != io.EOF {
			return oopsBuilder.
				With("bytesDownloaded", bytesDownloaded).
				With("bytesBuffer", bytesBuffer).
				Wrapf(err, "failed to read response body of: %s", url)
		}

		if bytesDownloaded == 0 {
			break
		}

		// Send the number of bytes downloaded to the channel
		bytesDownloadedChan <- bytesDownloaded

		// Write the bytes to the file
		if bytesWritten, err = downloadFile.Write(bytesBuffer[:bytesDownloaded]); err != nil {
			return oopsBuilder.
				With("bytesWritten", bytesWritten).
				With("bytesBuffer", bytesBuffer).
				With("bytesDownloaded", bytesDownloaded).
				Wrapf(err, "failed to write to file: %s", downloadFile.Name())
		}
	}

	return
}
