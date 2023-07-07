package utils

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AppendZipBar(p *mpb.Progress, download *DownloadResult) {

	// Config the bar
	zipBar := p.AddBar(0,
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(download.Dependency.Name),
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
	download.Dependency.ZipBar = zipBar

}
