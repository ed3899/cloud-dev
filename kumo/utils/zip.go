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

func Unzip(dr *DownloadResult, binsChan chan<- *Binary) {
	// 1. Open the zip file
	reader, err := zip.OpenReader(dr.Dependency.ZipPath)
	if err != nil {
		error := errors.Wrap(err, "failed to open zip file")
		binsChan <- &Binary{
			Dependency: dr.Dependency,
			Err:        error,
		}
		return
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err := filepath.Abs(dr.Dependency.ExtractionPath)
	if err != nil {
		error := errors.Wrap(err, "failed to get absolute path")
		binsChan <- &Binary{
			Dependency: dr.Dependency,
			Err:        error,
		}
		return
	}

	// 3. Iterate over zip files inside the archive and unzip each of them

	bytesUnzipped := make(chan int)
	unsuccesfulUnzip := make(chan bool, 1)

	// Wait group for unzipping goroutines
	var wgUnzip sync.WaitGroup

	// Wait group for unzipping goroutines
	for _, f := range reader.File {
		wgUnzip.Add(1)
		go func(f *zip.File) {
			defer wgUnzip.Done()

			bytesCopied, err := unzipFile(f, destination)
			if err != nil {
				error := errors.Wrap(err, "failed to unzip file")
				binsChan <- &Binary{
					Dependency: dr.Dependency,
					Err:        error,
				}
				unsuccesfulUnzip <- true
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(f)
	}

	go func() {
		wgUnzip.Wait()
		close(bytesUnzipped)
	}()

	// Update the progress bar for every unzipped file
	go func(dr *DownloadResult) {
		for b := range bytesUnzipped {
			dr.Dependency.ZipBar.IncrBy(b)
		}

		select {
			// If the unzipping was not successful, return
		case <-unsuccesfulUnzip:
			return
		default:
			// Otherwise, send the dependency to the channel
			binsChan <- &Binary{
				Dependency: dr.Dependency,
				Err:        nil,
			}
		}
	}(dr)

	wgUnzip.Wait()
}

func unzipFile(f *zip.File, destination string) (int64, error) {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return 0, fmt.Errorf("%s: illegal file path", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		return 0, nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return 0, err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return 0, err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer zippedFile.Close()

	bytesCopied, err := io.Copy(destinationFile, zippedFile)
	if err != nil {
		return 0, err
	}

	return bytesCopied, nil
}
