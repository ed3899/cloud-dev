package zip

import (
	"github.com/vbauerster/mpb/v8"
)

type ProgressBar interface {
	IncrBy(n int)
}

type MultiProgressBar interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}
