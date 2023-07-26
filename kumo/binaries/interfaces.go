package binaries

import "github.com/vbauerster/mpb/v8"

type Removable interface {
	Remove() error
}

type Retrivable interface {
	GetName() string
	GetPath() string
}

type Downloadable interface {
	SetDownloadBar(*mpb.Progress)
	Download(chan<- int) error
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(*mpb.Progress, int64)
	ExtractTo(string, chan<- int) error
	IncrementExtractionBar(int)
}

type ZipI interface {
	Retrivable
	Downloadable
	Extractable
	Removable
}