package download

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func BarSetDownloadingWith(
	progress IAddBar,
) ForIDownloader {
	return func(d IDownloader) {
		d.Bar().SetDownloading(
			progress.AddBar(int64(d.ContentLength()),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name(d.Name()),
					decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.Percentage(decor.WCSyncSpace),
						"downloaded",
					),
				),
			),
		)
	}
}

type IDownloader interface {
	IBarGetter
	IContentLengthGetter
	INameGetter
}

type ForIDownloader func(IDownloader)
