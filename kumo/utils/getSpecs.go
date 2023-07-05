package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"text/template"

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

func init() {
	hs := getHostSpecs()
	packerUrl := getPackerUrl(hs)
	downloadPacker(*packerUrl)
	// downloadPackages(hs)
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

type PackerURL struct {
	URL string
}

func getPackerUrl(s Specs) *PackerURL {
	packerTemplateUrl := "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_{{ .OS }}_{{ .ARCH }}.zip"
	tmpl, err := template.New("packer-url").Parse(packerTemplateUrl)

	if err != nil {
		log.Fatalf("error parsing template: %s. Here is the error: %#v\n", tmpl.Name(), err)
	}

	var packerTemplate bytes.Buffer
	err = tmpl.Execute(&packerTemplate, s)
	if err != nil {
		log.Fatalf("error executing template: %#v", err)
	}

	return &PackerURL{
		URL: packerTemplate.String(),
	}

}

func downloadPacker(p PackerURL) {
	downloadFile(p.URL, "./packer.zip")
}

func downloadFile(url string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("there was an error while attempting to download %#v", url)
	}
	defer response.Body.Close()

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		log.Fatalf("there was an error while creating %#v", filePath)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		response.ContentLength,
		filePath,
	)

	_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	if err != nil {
		log.Fatalf("there was an error while copying contents to %#v", filePath)
		return err
	}

	return nil
}

func downloadPackages(s Specs) {

}
