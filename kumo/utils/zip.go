package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Unzip(dr *DownloadResult, binsChan chan *Binary) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(dr.Dependency.ZipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err := filepath.Abs(dr.Dependency.ExtractionPath)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them

	bytesUnzipped := make(chan int, 4096)
	// bins := make(chan *Binary, 8)
	errCh := make(chan error, len(reader.File))

	// Wait group for unzipping goroutines
	var wgUnzip sync.WaitGroup
	wgUnzip.Add(len(reader.File))

	// Wait group for unzipping goroutines
	for _, f := range reader.File {
		go func(f *zip.File) {
			defer wgUnzip.Done()

			bytesCopied, err := unzipFile(f, destination)
			if err != nil {
				log.Printf("there was an error while unzipping the file: %v", err)
				errCh <- err
				binsChan <- &Binary{
					Dependency: dr.Dependency,
					Extracted:  false,
					Err:        err,
				}
				close(binsChan)
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(f)
	}

	// Wait for all unzipping goroutines to finish
	go func() {
		wgUnzip.Wait()
		close(errCh)
		// close(bins)
	}()

	// Update the progress bar for every unzipped file
	go func(dr *DownloadResult) {
		for b := range bytesUnzipped {
			dr.Dependency.ZipBar.IncrBy(b)
		}
		close(bytesUnzipped)

		binsChan <- &Binary{
			Dependency: dr.Dependency,
			Extracted:  true,
			Err:        nil,
		}
		close(binsChan)
	}(dr)

	// Range over the error channel and return the first error
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
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
