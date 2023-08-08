package interfaces

import (
	"github.com/vbauerster/mpb/v8"
)

type MpbV8BarIncrementor interface {
	IncrBy(n int)
}

type MpbV8MultiprogressBar interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}
