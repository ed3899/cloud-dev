package download

// Shuts down the progress bar.
func (d *Download) ProgressShutdown() {
	d.Progress.Shutdown()
}
