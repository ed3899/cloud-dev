package utils

type Specs struct {
	OS   string
	ARCH string
}

type Dependency struct {
	URL            string
	ZipPath        string
	ExtractionPath string
	Present        bool
}

type DownloadResult struct {
	Dependency *Dependency
	Err        error
}
