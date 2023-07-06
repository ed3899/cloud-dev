package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"path/filepath"

	"github.com/vbauerster/mpb/v8"
)

func DownloadPackages(ze []*ZipExecutableRef) {
	resultChan := make(chan DownloadResult, len(ze))
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))
	bars := GenerateBars(progress, ze)
	wg.Add(len(ze))

	for i, v := range ze {
		go func(ref *ZipExecutableRef, bar *mpb.Bar) {
			defer wg.Done()
			DownloadZip(ref, resultChan, bar)
		}(v, bars[i])
	}

	go func() {
		wg.Wait()

		close(resultChan)
	}()

	for result := range resultChan {
		if result.Err != nil {
			log.Printf("Error downloading %s: %v\n", result.ZipRef.URL, result.Err)
		}
	}
}

func DownloadZip(ref *ZipExecutableRef, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := Download(ref.URL, ref.BinPath, bar)

	result := DownloadResult{
		ZipRef: ref,
		Err:    err,
	}

	resultChan <- result
}

func Download(url string, binPath string, bar *mpb.Bar) error {
	// Download
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}
	defer response.Body.Close()

	buffer := make([]byte, 4096)

	for {
		bytesDownloaded, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("err while downloading from '%s': %#v", url, err)
		}
		if bytesDownloaded == 0 {
			break // Reached the end of the response body
		}
		bar.IncrBy(bytesDownloaded)
	}

	// Create the file along with all the necessary directories
	err = os.MkdirAll(filepath.Dir(binPath), 0755)
	if err != nil {
		log.Printf("there was an error while creating %#v", binPath)
		return err
	}

	// Create
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Printf("there was an error while creating %#v", binPath)
		return err
	}
	defer file.Close()

	// Fill
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Printf("there was an error while copying contents to %#v", binPath)
		return err
	}

	return nil
}