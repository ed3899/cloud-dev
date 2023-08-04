package interfaces

type Removable interface {
	Remove() error
}

type Retrivable interface {
	GetName() string
	GetPath() string
}

type Downloadable interface {
	SetDownloadBar(MultiProgressBar)
	Download(chan<- int) error
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(MultiProgressBar, int64)
	ExtractTo(string, chan<- int) error
	IncrementExtractionBar(int)
}

type Zip interface {
	Retrivable
	Downloadable
	Extractable
	Removable
}
