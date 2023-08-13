package download

func IncrementDownloading(
	downloadedBytes int,
) ForIBarGetter {
	barGetter := func(d IBarGetter) {
		d.Bar().Downloading().IncrBy(downloadedBytes)
	}

	return barGetter
}

type ForIBarGetter func(IBarGetter)
