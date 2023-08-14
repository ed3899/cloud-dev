package download

import "github.com/ed3899/kumo/common/interfaces"

func DownloadBarExtractingIncrByWith(
	extractedBytes int,
) DownloadBarExtractingIncryBy[Download] {
	downloadBarExtractingIncrBy := func(download interfaces.IClone[Download]) Download {
		downloadClone := download.Clone()

		downloadClone.Bar().Clone().Extracting().IncrBy(extractedBytes)

		return downloadClone
	}

	return downloadBarExtractingIncrBy
}

type DownloadBarExtractingIncryBy[D IDownload] func(interfaces.IClone[D]) D
