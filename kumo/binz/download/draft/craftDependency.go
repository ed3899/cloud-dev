package draft

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ed3899/kumo/host"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type Dependency struct {
	Name           string
	Present        bool
	URL            string
	ZipPath        string
	ExtractionPath string
	ContentLength  int64
	DownloadBar    *mpb.Bar
	ZipBar         *mpb.Bar
}

// Craft a dependency. If the dependency is not present, it will be downloaded.
//
// Given that in your directory structure, you have already downloaded the dependency
// "packer" and the executable is present in the directory `{cwd}\{dependenciesDirName}\packer\packer.exe`.
//
//	Example:
//
//	dependency, err := CraftHashicorpDependency("packer")
//
//	Result:
//
//	dependency = &Dependency{
//		Name:           "packer",
//		Present:        true,
//		URL:            "https://releases.hashicorp.com/packer/{version}/packer_{version}_{os}_{arch}.zip",
//		ExtractionPath: "{cwd}/{dependenciesDirName}/packer",
//		ZipPath:        "{cwd}/{dependenciesDirName}/packer/packer_{version}_{os}_{arch}.zip",
//		ContentLength:  {int},
//	}
//
// The paths are absolute paths and the format changes depending on the OS.
func CraftHashicorpDependency(name string) (dp *Dependency, err error) {
	depsDirName := utils.GetDependenciesDirName()
	specs := host.GetSpecs()
	zipName := fmt.Sprintf("%s_%s_%s.zip", name, specs.OS, specs.ARCH)

	destinationZipPath, err := filepath.Abs(filepath.Join(depsDirName, zipName))
	if err != nil {
		msg := fmt.Sprintf("failed to get zip path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	destinationExtractionPath, err := filepath.Abs(filepath.Join(depsDirName, name))
	if err != nil {
		msg := fmt.Sprintf("failed to get extraction path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	url, err := utils.GetDependencyURL(name, specs)
	if err != nil {
		msg := fmt.Sprintf("failed to get url for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		msg := fmt.Sprintf("failed to get content length for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	executablePath, err := filepath.Abs(filepath.Join(depsDirName, name, fmt.Sprintf("%s.exe", name)))
	if err != nil {
		msg := fmt.Sprintf("failed to get executable path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	if utils.FileNotPresent(executablePath) {
		log.Printf("%s not present", name)
		return &Dependency{
			Name:           name,
			Present:        false,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			ContentLength:  contentLength,
		}, nil
	}

	return &Dependency{
		Name:           name,
		Present:        true,
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		ContentLength:  contentLength,
	}, nil
}
