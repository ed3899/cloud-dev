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

	"github.com/vbauerster/mpb/v8"
	// "github.com/vbauerster/mpb/v8/decor"
)

func UnzipSource(dr *DownloadResult, wg *sync.WaitGroup) *Binary {
	// 1. Open the zip file
	defer wg.Done()

	reader, err := zip.OpenReader(dr.Dependency.ZipPath)
	if err != nil {
		log.Printf("there was an error while opening the zip file: %v", err)
		return &Binary{
			Dependency: dr.Dependency,
			Extracted:  false,
			Err:        err,
		}
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err := filepath.Abs(dr.Dependency.ExtractionPath)
	if err != nil {
		log.Printf("there was an error while getting the absolute path: %v", err)
		return &Binary{
			Dependency: dr.Dependency,
			Extracted:  false,
			Err:        err,
		}
	}

	// // Create a progress bar for unzipping the files
	// progress := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))
	// unzipBar := progress.AddBar(int64(len(reader.File)),
	// 	mpb.BarQueueAfter(dr.Dependency.DownloadBar),
	// 	mpb.BarFillerClearOnComplete(),
	// 	mpb.PrependDecorators(
	// 		decor.Name(fmt.Sprintf("%s.zip", dr.Dependency.Name)),
	// 		decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
	// 	),
	// 	mpb.AppendDecorators(
	// 		decor.OnComplete(
	// 			decor.Percentage(decor.WCSyncSpace),
	// 			"done",
	// 		),
	// 	),
	// )

	// for i := 0; i < len(reader.File); i++ {
	// 	unzipBar.IncrBy(i)
	// }

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination, dr.Dependency.ZipBar)
		if err != nil {
			log.Printf("there was an error while unzipping the file: %v", err)
			return &Binary{
				Dependency: dr.Dependency,
				Extracted:  false,
				Err:        err,
			}
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string, bar *mpb.Bar) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
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
		bar.IncrBy(bytesCopied)
	}

	return nil
}
