package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/schollz/progressbar/v3"
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

// func downloadPacker(p ZipExecutableRef) {
// 	downloadBin(p.URL, p.BinPath)
// }

func downloadPackages(ze []*ZipExecutableRef) {
	resultChan := make(chan DownloadResult, 2)

	for i := 0; i < len(ze); i++ {
		go downloadBin(*ze[i], resultChan)
	}

	for i := 0; i < len(ze); i++ {
		result := <-resultChan
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

func downloadBin(ref ZipExecutableRef, resultChan chan<- DownloadResult) {
	err := download(ref.URL, ref.BinPath)

	result := DownloadResult{
		ZipRef: &ref,
		Err:    err,
	}

	resultChan <- result
}

func download(url string, binPath string) error {
	// Download
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("there was an error while attempting to download %#v", url)
	}
	defer response.Body.Close()

	// Create
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Fatalf("there was an error while creating %#v", binPath)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		response.ContentLength,
		fmt.Sprintf("Downloading %s", binPath),
	)

	// Fill
	_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	if err != nil {
		log.Fatalf("there was an error while copying contents to %#v", binPath)
		return err
	}

	return nil
}
