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

func UnzipSource(dr *DownloadResult, wg *sync.WaitGroup) error {
	// 1. Open the zip file
	defer wg.Done()

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
	bins := make(chan *Binary, 8)
	errCh := make(chan error, len(reader.File))

	var wgUnzip sync.WaitGroup // Wait group for unzipping goroutines
	wgUnzip.Add(len(reader.File))

	for _, f := range reader.File {
		go func(f *zip.File) {
			defer wgUnzip.Done()

			bytesCopied, err := unzipFile(f, destination)
			if err != nil {
				log.Printf("there was an error while unzipping the file: %v", err)
				errCh <- err
				bins <- &Binary{
					Dependency: dr.Dependency,
					Extracted:  false,
					Err:        err,
				}
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(f)
	}

	go func() {
		wgUnzip.Wait()
		close(errCh)
		close(bins)
	}()

	go func() {
		for b := range bytesUnzipped {
			dr.Dependency.ZipBar.IncrBy(b)
		}
		close(bytesUnzipped)
	}()

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

	// buffer := make([]byte, 4096)
	// // 8. Copy the content of the file to the destination file while updating the progress bar
	// for {
	// 	bytesCopied, err := zippedFile.Read(buffer)
	// 	if err != nil && err != io.EOF {
	// 		return err

	// 	}
	// 	if bytesCopied == 0 {
	// 		break
	// 	}
	// 	if _, err := destinationFile.Write(buffer[:bytesCopied]); err != nil {
	// 		return err
	// 	}

	// 	bytesUnzipped <- bytesCopied
	// }

	// Give me an alternative using io.Copy

	bytesCopied, err := io.Copy(destinationFile, zippedFile)
	if err != nil {
		return 0, err
	}

	return bytesCopied, nil
}
