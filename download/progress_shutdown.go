package download

func ProgressShutdown(
	download *Download,
) {
	download.Progress.Shutdown()
}
