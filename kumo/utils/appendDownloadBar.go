package utils

import (
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AppendDownloadBar(wg *sync.WaitGroup, deps []*Dependency) *mpb.Progress {
	progress := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))

	for _, d := range deps {
		// Config the bar
		downloadBar := progress.AddBar(int64(d.ContentLength),
			mpb.BarFillerClearOnComplete(),
			mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),
			),
		)

		// Assign the bar to the dependency
		d.DownloadBar = downloadBar
	}

	return progress
}
