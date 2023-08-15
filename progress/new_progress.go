package progress

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/download"
	"github.com/ed3899/kumo/utils/url"
	"github.com/ed3899/kumo/utils/zip"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func NewProgress(
	options ...Option,
) (
	progress *Progress,
) {
	var (
		opt Option
	)

	progress = &Progress{}
	for _, opt = range options {
		opt(progress)
	}

	return
}

func WithDownload(d *download.Download) (option Option) {
	option = func(progress *Progress) {
		progress.Download = d

		return
	}

	return
}

func WithMultiProgressBar(p interfaces.MpbV8MultiProgressBar) (option Option) {
	pr := mpb.WithWaitGroup()
	option = func(progress *Progress) {
		progress.MultiProgressBar = p

		return
	}

	return
}

type Option func(*Progress)

type Progress struct {
	Download         *download.Download
	MultiProgressBar interfaces.MpbV8MultiProgressBar
}

func (dae *Progress) DownloadAndShowProgress(
	multiProgressBar interfaces.MpbV8MultiProgressBar,
	urlDownload url.DownloadF,
) (err error) {
	var (
		downloadedBytesChan = make(chan int, 1024)
		errChan             = make(chan error, 1)
		doneChan            = make(chan bool, 1)
		oopsBuilder         = oops.
					Code("download_and_show_progress_failed").
					With("download.Name", dae.Download.Name).
					With("download.AbsPath", dae.Download.AbsPath).
					With("multiProgressBar", multiProgressBar)

		downloadedBytes int
		done            bool
	)

	go func() {
		defer close(downloadedBytesChan)
		defer close(errChan)
		defer close(doneChan)

		dae.Download.SetDownloadBar(multiProgressBar)

		if err = urlDownload(dae.Download.Url, dae.Download.AbsPath, downloadedBytesChan); err != nil {
			err = oopsBuilder.
				With("url", dae.Download.Url).
				With("absPath", dae.Download.AbsPath).
				With("downloadedBytesChan", downloadedBytesChan).
				Wrapf(err, "failed to download: %v", dae.Download.Url)
			errChan <- err
			return
		}

		doneChan <- true
	}()

OuterLoop:
	for {
		select {
		case downloadedBytes = <-downloadedBytesChan:
			if downloadedBytes > 0 {
				dae.Download.IncrementDownloadBar(downloadedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while downloading %s", dae.Download.Name)
				return
			}

		case done = <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return
}

func (dae *Progress) ExtractAndShowProgress(
	multiProgressBar interfaces.MpbV8MultiProgressBar,
	getZipSize zip.GetZipSizeF,
	unzip zip.UnzipF,
) (err error) {
	var (
		extractedBytesChan = make(chan int, 1024)
		errChan            = make(chan error, 1)
		doneChan           = make(chan bool, 1)
		oopsBuilder        = oops.
					Code("ExtractAndShowProgress").
					With("download.Name", dae.Download.Name).
					With("download.AbsPath", dae.Download.AbsPath).
					With("multiProgressBar", multiProgressBar)

		extractedBytes int
		done           bool
		zipSize        int64
	)

	if zipSize, err = getZipSize(dae.Download.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip size for: %v", dae.Download.AbsPath)
		return
	}

	go func(zipSize int64) {
		defer close(errChan)
		defer close(doneChan)

		dae.Download.SetExtractionBar(multiProgressBar, zipSize)

		if err = unzip(dae.Download.AbsPath, filepath.Dir(dae.Download.AbsPath), extractedBytesChan); err != nil {
			err = oopsBuilder.
				With("absPath", dae.Download.AbsPath).
				With("extractedBytesChan", extractedBytesChan).
				Wrapf(err, "failed to unzip: %v", dae.Download.AbsPath)
			errChan <- err
			return
		}

		doneChan <- true
	}(zipSize)

OuterLoop:
	for {
		select {
		case extractedBytes = <-extractedBytesChan:
			if extractedBytes > 0 {
				dae.Download.IncrementExtractionBar(extractedBytes)
			}

		case err = <-errChan:
			if err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while extracting %s", dae.Download.Name)
				return
			}

		case done = <-doneChan:
			if done {
				break OuterLoop
			}
		}
	}

	return
}

func (dae *Progress) CleanAbsPathToZipDirWith(
	removeAll RemoverF,
) (err error) {
	var (
		oopsBuilder = oops.
				Code("RemoveAllWith").
				With("removeAll", removeAll)
		absPathToZipDir = filepath.Dir(dae.Download.AbsPath)
	)

	if err = removeAll(absPathToZipDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return
	}

	return
}

func (dae *Progress) RemoveDownloadWith(
	remove RemoverF,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("RemoveWith").
			With("remove", remove)
	)

	if err = remove(dae.Download.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", dae.Download.AbsPath)
		return
	}

	return
}

func (dae *Progress) Shutdown() {
	dae.MultiProgressBar.Shutdown()
}

type RemoverF func(path string) (err error)
