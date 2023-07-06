package utils

import "github.com/vbauerster/mpb/v8"

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
	Bar            *mpb.Bar
}

type DownloadResult struct {
	Dependency *Dependency
	Err        error
}

type Binary struct {
	Dependency *Dependency
	Extracted	bool
	Err error
}
