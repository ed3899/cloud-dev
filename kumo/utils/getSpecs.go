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
	"time"

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
	validateHostCompatibility(hs)
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

func generateBars(progress *mpb.Progress, ze []*ZipExecutableRef) []*mpb.Bar {
	bars := make([]*mpb.Bar, 0)

	for i := 0; i < len(ze); i++ {
		var bar *mpb.Bar
		url := ze[i].URL
		name := ze[i].BinPath
		resp, err := http.Head(url)
		if err != nil {
			log.Printf("Error occurred while sending HEAD request: %v\n", err)
			continue
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-200 status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()

		// Retrieve the Content-Length header
		contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			log.Printf("There was an error while attempting to parse the content length from '%s': '%#v'", url, err)
			resp.Body.Close()
			continue
		}

		bar = progress.AddBar(int64(contentLength),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					// ETA decorator with ewma age of 30
					decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), "done",
				),
			),
		)

		bars = append(bars, bar)
	}

	return bars
}

func downloadPackages(ze []*ZipExecutableRef) {
	resultChan := make(chan DownloadResult)
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))
	bars := generateBars(progress, ze)
	wg.Add(len(ze))

	for i := 0; i < len(ze); i++ {
		go func(ref ZipExecutableRef, bar *mpb.Bar) {
			defer wg.Done()
			downloadBin(ref, resultChan, bar)
		}(*ze[i], bars[i])
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

type DownloadResult struct {
	ZipRef *ZipExecutableRef
	Err    error
}

func downloadBin(ref ZipExecutableRef, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := download(ref.URL, ref.BinPath, bar)

	result := DownloadResult{
		ZipRef: &ref,
		Err:    err,
	}

	resultChan <- result
}

func download(url string, binPath string, bar *mpb.Bar) error {
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
