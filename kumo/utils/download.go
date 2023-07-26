package utils

import (
	"io"
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
			With("err", err).
			Errorf("failed to download from: %s", url)
	}
	defer func() (err error) {
		if err = response.Body.Close(); err != nil {
			return oopsBuilder.
				With("err", err).
				With("response", response).
				Errorf("failed to close response body")
		}
		return
	}()

	// Create the destination dir
	if err = os.MkdirAll(destDir, 0755); err != nil {
		return oopsBuilder.
			With("err", err).
			Errorf("failed to create destination directory for: %s", destPath)
	}

	// Create the file to write to
	if downloadFile, err = os.Create(destPath); err != nil {
		return oopsBuilder.
			With("err", err).
			Errorf("failed to create file for: %s", destPath)
	}
	defer func() (err error) {
		if err = downloadFile.Close(); err != nil {
			return oopsBuilder.
				With("err", err).
				With("file", downloadFile).
				Errorf("failed to close file")
		}
		return
	}()

	// Iterate over the response body
	for {
		// Read the response body into the bytes buffer
		if bytesDownloaded, err = response.Body.Read(bytesBuffer); err != nil && err != io.EOF {
			return oopsBuilder.
				With("err", err).
				With("bytesDownloaded", bytesDownloaded).
				With("bytesBuffer", bytesBuffer).
				Errorf("failed to read response body")
		}

		if bytesDownloaded == 0 {
			break
		}

		// Send the number of bytes downloaded to the channel
		bytesDownloadedChan <- bytesDownloaded

		// Write the bytes to the file
		if bytesWritten, err = downloadFile.Write(bytesBuffer[:bytesDownloaded]); err != nil {
			return oopsBuilder.
				With("err", err).
				With("bytesWritten", bytesWritten).
				With("bytesBuffer", bytesBuffer).
				With("bytesDownloaded", bytesDownloaded).
				Errorf("failed to write to file")
		}
	}

	return
}
