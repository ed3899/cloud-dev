package utils

import (
	// "sync"
	// "time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AttachDownloadBar(p *mpb.Progress, d *Dependency) {
	// Config the bar
	downloadBar := p.AddBar(int64(d.ContentLength),
		mpb.BarFillerClearOnComplete(),
		// mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.Name(d.Name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"downloaded",
			),
		),
	)

	// Assign the bar to the dependency
	d.DownloadBar = downloadBar
}
