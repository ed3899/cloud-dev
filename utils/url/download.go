package url

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

func Download(url, path string, bytesDownloadedChan chan<- int) error {
	oopsBuilder := oops.
		Code("DownloadWith").
		In("utils").
		In("url").
		Tags("Download")

	// Initiate download and defer closing the response body
	response, err := http.Get(url)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to initiate download from: %s", url)
		return err
	}
	defer response.Body.Close()

	// Create the destination dir
	destDir := filepath.Dir(path)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to create destination directory for: %s", destDir)
		return err
	}

	// Create the file to write to and defer closing it
	downloadFile, err := os.Create(path)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to create file for: %s", path)
		return err
	}
	defer downloadFile.Close()

	// Iterate over the response body
	bytesBuffer := make([]byte, 4096)
	for {
		// Read the response body into the bytes buffer
		bytesDownloaded, err := response.Body.Read(bytesBuffer)
		if err != nil && err != io.EOF {
			err := oopsBuilder.
				With("bytesDownloaded", bytesDownloaded).
				With("bytesBuffer", bytesBuffer).
				Wrapf(err, "failed to read response body of: %s", url)
			return err
		}

		if bytesDownloaded == 0 {
			break
		}

		// Send the number of bytes downloaded to the channel
		bytesDownloadedChan <- bytesDownloaded

		// Write the bytes to the file
		bytesWritten, err := downloadFile.Write(bytesBuffer[:bytesDownloaded])
		if err != nil {
			err := oopsBuilder.
				With("bytesWritten", bytesWritten).
				With("bytesBuffer", bytesBuffer).
				With("bytesDownloaded", bytesDownloaded).
				Wrapf(err, "failed to write to file: %s", downloadFile.Name())

			return err
		}
	}

	return nil
}
