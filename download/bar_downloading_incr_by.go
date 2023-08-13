package download

func BarDownloadingIncrBy(
	downloadedBytes int,
	download IBarGetter,
) {
	download.Bar().Downloading().IncrBy(downloadedBytes)
}
