package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Specs struct {
	OS   string
	ARCH string
}

type ZipExecutableRef struct {
	URL     string
	ZipPath string
	BinPath string
}

func init() {
	hs := getHostSpecs()
	vh := validateHostCompatibility(hs)
	packerUrl := getPackerUrl(vh)
	pulumiUrl := getPulumiUrl(vh)
	urls := []*ZipExecutableRef{packerUrl, pulumiUrl}
	downloadPackages(urls)
}

func getHostSpecs() Specs {
	return Specs{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}

func validateHostCompatibility(s Specs) Specs {
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

	return s
}

func getPackerUrl(s Specs) *ZipExecutableRef {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("there was an error getting the current directory: %v", err)
	}

	// Create the destination path with dir + packer + packer.exe
	destinationZipPath := filepath.Join(dir, "packer", fmt.Sprintf("packer_%s_%s.zip", s.OS, s.ARCH))

	// Return the zip executable reference
	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_%s_%s.zip", s.OS, s.ARCH),
		BinPath: destinationZipPath,
	}
}

func getPulumiUrl(s Specs) *ZipExecutableRef {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("there was an error getting the current directory: %v", err)
	}

	// Create the destination path
	var arch string
	switch s.ARCH {
	case "amd64":
		arch = "x64"
	}
	destinationZipPath := filepath.Join(dir, "pulumi", fmt.Sprintf("packer_%s_%s.zip", s.OS, arch))

	// Return the zip executable reference
	return &ZipExecutableRef{
		URL:     fmt.Sprintf("https://get.pulumi.com/releases/sdk/pulumi-v3.74.0-%s-%s.zip", s.OS, arch),
		BinPath: destinationZipPath,
	}
}

func generateBars(progress *mpb.Progress, ze []*ZipExecutableRef) []*mpb.Bar {
	// Create the bars
	bars := make([]*mpb.Bar, 0)
	for i := 0; i < len(ze); i++ {
		var bar *mpb.Bar
		url := ze[i].URL
		name := filepath.Base(ze[i].BinPath)
		// Perform a HEAD request to get the content length
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

		// Assign the bar length to the content length
		bar = progress.AddBar(int64(contentLength),
			mpb.PrependDecorators(
				decor.Name(name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),
			),
		)

		bars = append(bars, bar)
	}

	return bars
}

// Write a function that checks for the existance of the binaries
func checkBinariesExistence(ze []*ZipExecutableRef) bool {
	for _, e := range ze {
		if _, err := os.Stat(e.BinPath); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Write a function that unzip the files
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the directory
	os.MkdirAll(dest, 0755)

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Open the file
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the file
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(fpath), f.Mode())
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			// Write the file
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func downloadPackages(ze []*ZipExecutableRef) {
	resultChan := make(chan DownloadResult, len(ze))
	wg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))
	bars := generateBars(progress, ze)
	wg.Add(len(ze))

	for i, v := range ze {
		go func(ref *ZipExecutableRef, bar *mpb.Bar) {
			defer wg.Done()
			downloadBin(ref, resultChan, bar)
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

type DownloadResult struct {
	ZipRef *ZipExecutableRef
	Err    error
}

func downloadBin(ref *ZipExecutableRef, resultChan chan<- DownloadResult, bar *mpb.Bar) {
	err := download(ref.URL, ref.BinPath, bar)

	result := DownloadResult{
		ZipRef: ref,
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
