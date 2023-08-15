package url

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

type IDownload interface {
	Url() string
	Path() string
}

func DownloadWith(
	httpGet func(url string) (resp *http.Response, err error),
	osMkdirAll func(path string, perm fs.FileMode) error,
	osCreate func(name string) (*os.File, error),
) Download {
	oopsBuilder := oops.
		Code("DownloadWith").
		In("utils").
		In("url").
		Tags("Download")

	download := func(download IDownload, bytesDownloadedChan chan<- int) error {
		// Initiate download and defer closing the response body
		response, err := http.Get(download.Url())
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to initiate download from: %s", download.Url())
			return err
		}
		defer response.Body.Close()

		// Create the destination dir
		destDir := filepath.Dir(download.Path())
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to create destination directory for: %s", destDir)
			return err
		}

		// Create the file to write to and defer closing it
		downloadFile, err := os.Create(download.Path())
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to create file for: %s", download.Path())
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
					Wrapf(err, "failed to read response body of: %s", download.Url())
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

	return download
}

type Download func(download IDownload, bytesDownloadedChan chan<- int) error
