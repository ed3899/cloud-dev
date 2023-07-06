package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	// "github.com/schollz/progressbar/v3"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Specs struct {
	OS   string
	ARCH string
}

type Urls struct {
	Packer string
	Pulumi string
}

type ZipExecutableRef struct {
	URL     string
	BinPath string
}

func init() {
	hs := getHostSpecs()
	packerUrl := getPackerUrl(hs)
	pulumiUrl := getPulumiUrl(hs)
	urls := []*ZipExecutableRef{packerUrl, pulumiUrl}

	// downloadPacker(*packerUrl)
	downloadPackages(urls)
}

func getHostSpecs() Specs {
	return Specs{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}

func validateHostCompatibility(s Specs) {
	switch s.OS {
	case "windows":
		switch s.ARCH {
		case "386":
		case "amd64":
		default:
			log.Fatalf("Looks like your operative systems architecture is not supported :/")
		}
	default:
		log.Fatalf("Looks like your operative system is not supported :/")
	}
}

func getPackerUrl(s Specs) *ZipExecutableRef {
	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH),
		BinPath: fmt.Sprintf("packer_%s_%s.zip", s.OS, s.ARCH),
	}
}

func getPulumiUrl(s Specs) *ZipExecutableRef {
	var arch string
	switch s.ARCH {
	case "amd64":
		arch = "x64"
	}

	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://get.pulumi.com/releases/sdk/pulumi-v3.74.0-%s-%s.zip", s.OS, arch),
		BinPath: fmt.Sprintf("pulumi_%s_%s.zip", s.OS, s.ARCH),
	}
}

func downloadPackages(ze []*ZipExecutableRef) {
	resultChan := make(chan DownloadResult)
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg))

	for i := 0; i < len(ze); i++ {
		wg.Add(1)
		go func(ref ZipExecutableRef) {
			defer wg.Done()
			downloadBin(ref, resultChan, progress)
		}(*ze[i])
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.Err != nil {
			log.Printf("Error downloading %s: %v\n", result.ZipRef.URL, result.Err)
		} else {
			fmt.Printf("Downloaded %s to %s successfully\n", result.ZipRef.URL, result.ZipRef.BinPath)
		}
	}
}

type DownloadResult struct {
	ZipRef *ZipExecutableRef
	Err    error
}

func downloadBin(ref ZipExecutableRef, resultChan chan<- DownloadResult, progress *mpb.Progress) {
	err := download(ref.URL, ref.BinPath, progress)

	result := DownloadResult{
		ZipRef: &ref,
		Err:    err,
	}

	resultChan <- result
}

func download(url string, binPath string, progress *mpb.Progress) error {
	// Download
	response, err := http.Get(url)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}
	defer response.Body.Close()

	totalSize, err := strconv.ParseInt(response.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Printf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}

	buffer := make([]byte, 4096)

	bar := progress.AddBar(int64(totalSize),
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(binPath),
			// decor.DSyncWidth bit enables column width synchronization
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				// ETA decorator with ewma age of 30
				decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), "done",
			),
		))

	for {
		bytesDownloaded, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			// Handle the error accordingly (e.g., log, return, etc.)
			log.Printf("err while downloading from '%s': %#v", url, err)
		}

		if bytesDownloaded == 0 {
			break // Reached the end of the response body
		}
		bar.IncrBy(bytesDownloaded)

	}

	// Create
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Fatalf("there was an error while creating %#v", binPath)
	}
	defer file.Close()

	// Fill
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatalf("there was an error while copying contents to %#v", binPath)
		return err
	}

	return nil
}
