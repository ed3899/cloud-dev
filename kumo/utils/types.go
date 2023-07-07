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
	Present        bool
	ContentLength  int64
	DownloadBar    *mpb.Bar
	ZipBar				 *mpb.Bar
}

type DownloadResult struct {
	Dependency *Dependency
	Err        error
}

type ZipFile struct {
	Name string
	Path string
	Extracted bool
	Error error
}

type Binary struct {
	Dependency *Dependency
	Extracted  bool
	Err        error
}
