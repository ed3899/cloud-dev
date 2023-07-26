package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func Unzip(pathToZip, extractToPath string, bytesUnzipped chan<- int) (err error) {
	var (
		reader  *zip.ReadCloser
		errChan chan error
		zipFile *zip.File

		unzipGroup = new(sync.WaitGroup)
	)

	// Open the zip file
	if reader, err = zip.OpenReader(pathToZip); err != nil {
		return errors.Wrap(err, "failed to open zip file")
	}
	defer func() {
		if err = reader.Close(); err != nil {
			err = errors.Wrap(err, "failed to close zip file")
			return
		}
	}()

	// Declare error channel
	errChan = make(chan error, len(reader.File))

	// Unzip each file concurrently
	for _, zipFile = range reader.File {
		unzipGroup.Add(1)
		go func(zf *zip.File) {
			defer unzipGroup.Done()

			var (
				bytesCopied int64
			)

			if bytesCopied, err = unzipFile(zf, extractToPath); err != nil {
				errChan <- errors.Wrapf(err, "failed to unzip file: %s", zf.Name)
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(zipFile)
	}

	// Close channels when all goroutines are done
	go func() {
		unzipGroup.Wait()
		close(errChan)
		close(bytesUnzipped)
	}()

	// Check for errors
	for err = range errChan {
		if err != nil {
			return
		}
	}

	return
}

func unzipFile(f *zip.File, extractToPath string) (bytesCopied int64, err error) {
	var (
		filePath        string
		destinationFile *os.File
		zippedFile      io.ReadCloser
	)
	
	// Check if file paths are not vulnerable to Zip Slip
	filePath = filepath.Join(extractToPath, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(extractToPath)+string(os.PathSeparator)) {
		return 0, fmt.Errorf("%s: illegal file path", filePath)
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		return 0, nil
	}

	if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return 0, err
	}

	// Create a destination file for unzipped content
	if destinationFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
		return 0, err
	}
	defer func() {
		if errorClosingDestFile := destinationFile.Close(); errorClosingDestFile != nil {
			err = errors.Wrap(errorClosingDestFile, "Error occurred while closing destination file")
		}
	}()

	// Unzip the content of a file and copy it to the destination file
	if zippedFile, err = f.Open(); err != nil {
		return 0, err
	}
	defer func() {
		if errorClosingZippedFile := zippedFile.Close(); errorClosingZippedFile != nil {
			err = errors.Wrap(errorClosingZippedFile, "Error occurred while closing zipped file")
		}
	}()

	if bytesCopied, err = io.Copy(destinationFile, zippedFile); err != nil {
		return 0, err
	}

	return
}
