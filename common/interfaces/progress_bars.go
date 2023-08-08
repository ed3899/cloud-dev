package interfaces

import (
	"github.com/vbauerster/mpb/v8"
)

type ProgressBarIncrementor interface {
	IncrBy(n int)
}

type ProgressBarAdder interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}
