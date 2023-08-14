package download

func BarDownloadingIncrByWith(
	downloadedBytes int,
) BarDownloadingIncryBy {
	barDownloadingIncryBy := func(download IBarGetter) {
		download.Bar().Downloading().IncrBy(downloadedBytes)
	}

	return barDownloadingIncryBy
}

type Cloneable interface {

}

type BarDownloadingIncryBy func(IBarGetter)
