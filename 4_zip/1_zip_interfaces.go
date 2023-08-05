package zip

type Downloadable interface {
	SetDownloadBar(MultiProgressBar)
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(MultiProgressBar, int64)
	IncrementExtractionBar(int)
}

type ZipI interface {
	Downloadable
	Extractable
}
