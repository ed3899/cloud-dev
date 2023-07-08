package utils

import (
	"github.com/vbauerster/mpb/v8"
)

type Specs struct {
	OS   string
	ARCH string
}

type Dependency struct {
	Name           string
	URL            string
	ZipPath        string
	ExtractionPath string
	ContentLength  int64
	DownloadBar    *mpb.Bar
	ZipBar         *mpb.Bar
}

type Dependencies = []*Dependency

type DownloadResult struct {
	Dependency *Dependency
	Err        error
}

type Binary struct {
	Dependency *Dependency
	Err        error
}

type Binaries struct {
	Packer *Binary
	Pulumi *Binary
}

type Packer struct {
	ExecutablePath string
}

type Pulumi struct {
	ExecutablePath string
}