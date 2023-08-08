package interfaces

import (
	"github.com/vbauerster/mpb/v8"
)

type MpbV8BarIncrementor interface {
	IncrBy(n int)
}

type MpbV8MultiProgressBarAdder interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}

type MpbV8MultiProgressBarShutter interface {
	Shutdown()
}

type MpbV8MultiProgressBar interface {
	MpbV8MultiProgressBarAdder
	MpbV8MultiProgressBarShutter
}
