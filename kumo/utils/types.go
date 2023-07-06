package utils

type Specs struct {
	OS   string
	ARCH string
}

type ZipExecutableRef struct {
	URL     string
	ZipPath string
	BinPath string
}

type DownloadResult struct {
	ZipRef *ZipExecutableRef
	Err    error
}