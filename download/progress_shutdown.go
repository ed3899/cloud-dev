package download

func (d *Download) ProgressShutdown() {
	d.Progress.Shutdown()
}
