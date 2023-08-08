package interfaces

type DownloadBarSetter interface {
	SetDownloadBar(MpbV8MultiprogressBar)
}

type DownloadBarIncrementor interface {
	IncrementDownloadBar(int)
}

type Downloadable interface {
	DownloadBarSetter
	DownloadBarIncrementor
}

type ExtractBarSetter interface {
	SetExtractionBar(MpbV8MultiprogressBar, int64)
}

type ExtractBarIncrementor interface {
	IncrementExtractionBar(int)
}

type Extractable interface {
	ExtractBarSetter
	ExtractBarIncrementor
}

type ZipI interface {
	Downloadable
	Extractable
}
