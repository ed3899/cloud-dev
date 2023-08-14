package download

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func BarSetDownloadingWith(
	progress IAddBar,
) BarSetDownloading {
	barSetDownloading := func(download IDownloader) {
		download.Bar().SetDownloading(
			progress.AddBar(int64(download.ContentLength()),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name(download.Name()),
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

	return barSetDownloading
}

type BarSetDownloading func(IDownloader)

type IDownloader interface {
	IBarGetter[IDownloadingSetter]
	IContentLengthGetter
	INameGetter
}
