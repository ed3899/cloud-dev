package download

func BarDownloadingIncrByWith(
	downloadedBytes int,
) BarDownloadingIncryBy {
	barDownloadingIncryBy := func(download IBarGetter[IDownloadingGetter]) {
		download.Bar().Downloading().IncrBy(downloadedBytes)
	}

	return barDownloadingIncryBy
}

type BarDownloadingIncryBy func(IBarGetter[IDownloadingGetter])
