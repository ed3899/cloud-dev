package utils

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"text/template"
)

type Specs struct {
	OS   string
	ARCH string
}

func init() {
	hs := getHostSpecs()
	downloadPackages(hs)
}

func getHostSpecs() Specs {
	return Specs{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}

func downloadPackages(s Specs) {
	packerTemplateUrl := "https://releases.hashicorp.com/packer/1.9.1/packer_1.9.1_{{ .OS }}_{{ .ARCH }}.zip"
	tmpl, err := template.New("packer-url").Parse(packerTemplateUrl)
	if err != nil {
		fmt.Printf("error parsing template: %s. Here is the error: %#v\n", tmpl.Name(), err)
		os.Exit(1)
	}

	var packerTemplate bytes.Buffer
	err = tmpl.Execute(&packerTemplate, s)
	if err != nil {
		fmt.Println("error executing template:", err)
		os.Exit(1)
	}
	packerUrl := packerTemplate.String()

	switch s.OS {
	case "windows":
		switch s.ARCH {
		case "386":
			fmt.Println("Running windows 386")
		case "amd64":
			fmt.Println(packerUrl)
		default:
			panic("Looks like your operative systems architecture is not supported :/")
		}
	default:
		panic("Looks like your operative system is not supported :/")
	}
}
