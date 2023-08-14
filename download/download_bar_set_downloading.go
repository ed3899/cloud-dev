package download

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func DownloadBarSetDownloadingWith(
	progress IAddBar,
) DownloadBarSetDownloading[Download] {
	downloadBarSetDownloading := func(download interfaces.IClone[Download]) Download {
		downloadClone := download.Clone()

		downloadClone.Bar().SetDownloading(
			progress.AddBar(int64(downloadClone.ContentLength()),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name(downloadClone.Name()),
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

		return downloadClone
	}

	return downloadBarSetDownloading
}

type DownloadBarSetDownloading[D IDownload] func(interfaces.IClone[D]) D
