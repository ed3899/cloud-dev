package download

import "github.com/ed3899/kumo/common/interfaces"

func DownloadBarDownloadingIncrByWith(
	downloadedBytes int,
) DownloadBarDownloadingIncryBy[Download] {
	downloadBarDownloadingIncryBy := func(download interfaces.IClone[Download]) Download {
		downloadClone := download.Clone()

		downloadClone.Bar().Clone().Downloading().IncrBy(downloadedBytes)

		return downloadClone
	}

	return downloadBarDownloadingIncryBy
}

type DownloadBarDownloadingIncryBy[D IDownload] func(interfaces.IClone[D]) D
