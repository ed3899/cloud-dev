package download

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func DownloadBarSetExtractingWith(
	progress IAddBar,
	zipSize int64,
) DownloadBarSetExtracting[Download] {
	downloadBarSetExtracting := func(download interfaces.IClone[Download]) Download {
		downloadClone := download.Clone()

		downloadClone.Bar().SetExtracting(
			progress.AddBar(zipSize,
				mpb.BarQueueAfter(downloadClone.Bar().Downloading()),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name(downloadClone.Name()),
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

		return downloadClone
	}

	return downloadBarSetExtracting
}

type DownloadBarSetExtracting[D IDownload] func(interfaces.IClone[D]) D
