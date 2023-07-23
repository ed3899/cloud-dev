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
	// Open the zip file
	reader, err := zip.OpenReader(pathToZip)
	if err != nil {
		err = errors.Wrap(err, "failed to open zip file")
		return
	}
	defer reader.Close()

	// Wait group for unzipping goroutines
	unzipGroup := new(sync.WaitGroup)
	errChan := make(chan error, len(reader.File))

	// Unzip each file concurrently
	for _, f := range reader.File {
		unzipGroup.Add(1)
		go func(f *zip.File) {
			defer unzipGroup.Done()

			bytesCopied, err := unzipFile(f, extractToPath)
			if err != nil {
				err = errors.Wrapf(err, "failed to unzip file: %s", f.Name)
				errChan <- err
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(f)
	}

	// Close channels when all goroutines are done
	go func() {
		unzipGroup.Wait()
		close(errChan)
		close(bytesUnzipped)
	}()

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return
}

func unzipFile(f *zip.File, extractToPath string) (bytesCopied int64, err error) {
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(extractToPath, f.Name)
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
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return 0, err
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer zippedFile.Close()

	bytesCopied, err = io.Copy(destinationFile, zippedFile)
	if err != nil {
		return 0, err
	}

	return bytesCopied, nil
}
