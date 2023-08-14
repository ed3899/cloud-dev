package download

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func BarSetExtractingWith(
	progress IAddBar,
	zipSize int64,
) BarSetExtracting {
	barSetExtracting := func(extracter IExtracter) {
		extracter.Bar().SetExtracting(
			progress.AddBar(zipSize,
				mpb.BarQueueAfter(extracter.Bar().Downloading().(*mpb.Bar)),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name(extracter.Name()),
					decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.Percentage(decor.WCSyncSpace),
						"unzipped",
					),
				),
			),
		)
	}

	return barSetExtracting
}

type BarSetExtracting func(IExtracter)

type IExtracter interface {
	INameGetter
	IBarGetter[IDownloadingAndSetExtracting]
}

type IDownloadingAndSetExtracting interface {
	IExtractingSetter
	IDownloadingGetter
}
