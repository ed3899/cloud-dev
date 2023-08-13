package download

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func BarSetDownloading(
	progress IAddBar,
	download IDownloader,
) {

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

type IDownloader interface {
	IBarGetter
	IContentLengthGetter
	INameGetter
}
