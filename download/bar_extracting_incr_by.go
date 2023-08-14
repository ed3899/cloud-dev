package download

func BarExtractingIncrByWith(
	extractedBytes int,
) BarExtractingIncryBy {
	barExtractingIncrBy := func(download IBarGetter[IExtractingGetter]) {
		download.Bar().Extracting().IncrBy(extractedBytes)
	}

	return barExtractingIncrBy
}

type BarExtractingIncryBy func(IBarGetter[IExtractingGetter])
