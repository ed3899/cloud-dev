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
	// "time"
	// "github.com/vbauerster/mpb/v8/decor"
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
	zipFiles := make(chan *ZipFile, 8)
	done := make(chan bool, 5)

	var wgUnzip sync.WaitGroup // Wait group for unzipping goroutines
	wgUnzip.Add(len(reader.File))

	for _, f := range reader.File {
		go func(f *zip.File) {
			defer wgUnzip.Done()

			err := unzipFile(f, destination, bytesUnzipped)
			if err != nil {
				log.Printf("there was an error while unzipping the file: %v", err)
				bins <- &Binary{
					Dependency: dr.Dependency,
					Extracted:  false,
					Err:        err,
				}
			}
		}(f)
	}

	go func() {
		wgUnzip.Wait() // Wait for all unzipping goroutines to finish
		close(bins)
	}()

	go func() {
		for b := range bytesUnzipped {
			dr.Dependency.ZipBar.IncrBy(b)
		}
	}()

	wgUnzip.Wait()
	close(bytesUnzipped)

	go func() {
		for z := range zipFiles {
			if z.Error != nil {
				log.Printf("there was an error while unzipping the file: %v", z.Error)
				bins <- &Binary{
					Dependency: dr.Dependency,
					Extracted:  false,
					Err:        z.Error,
				}
			}
		}
		close(zipFiles)
		done <- true
	}()

	for b := range bins {
		if b.Err != nil {
			return b.Err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string, bytesUnzipped chan<- int) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	buffer := make([]byte, 4096)
	// 8. Copy the content of the file to the destination file while updating the progress bar
	for {
		bytesCopied, err := zippedFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err

		}
		if bytesCopied == 0 {
			break
		}
		if _, err := destinationFile.Write(buffer[:bytesCopied]); err != nil {
			return err
		}

		bytesUnzipped <- bytesCopied
		// bar.IncrBy(bytesCopied)
	}

	return nil
}
