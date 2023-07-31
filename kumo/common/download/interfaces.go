package download

import (
	"github.com/vbauerster/mpb/v8"
)

type ProgressBarI interface {
	IncrBy(n int)
}

type MultiProgressBarI interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}

type Removable interface {
	Remove() error
}

type Retrivable interface {
	GetName() string
	GetPath() string
}

type Downloadable interface {
	SetDownloadBar(MultiProgressBarI)
	Download(chan<- int) error
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(MultiProgressBarI, int64)
	ExtractTo(string, chan<- int) error
	IncrementExtractionBar(int)
}

type ZipI interface {
	Retrivable
	Downloadable
	Extractable
	Removable
}
