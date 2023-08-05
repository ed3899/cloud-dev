package zip

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/0_constants"

	"github.com/ed3899/kumo/common/zip/interfaces"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Zip struct {
	name          string
	absPath       string
	url           string
	contentLength int64
	downloadBar   *mpb.Bar
	extractionBar *mpb.Bar
}

func New(toolConfig common_tool_interfaces.Tool) (zip interfaces.Zip, err error) {
	var (
		oopsBuilder = oops.Code("zip_new_failed")
		absPath     = filepath.Join(dirs.DEPENDENCIES_DIR_NAME, toolConfig.Name(), fmt.Sprintf("%s.zip", toolConfig.Name()))

		contentLength int64
	)

	if contentLength, err = utils.GetContentLength(toolConfig.Url()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip content length")
		return
	}

	zip = &Zip{
		name:          filepath.Base(absPath),
		absPath:       absPath,
		url:           toolConfig.Url(),
		contentLength: contentLength,
	}

	return
}

func (z *Zip) Name() (name string) {
	return z.name
}

func (z *Zip) AbsPath() (path string) {
	return z.absPath
}

func (z *Zip) IsPresent() (present bool) {
	return utils.FilePresent(z.absPath)
}

func (z *Zip) IsNotPresent() (notPresent bool) {
	return utils.FileNotPresent(z.absPath)
}

func (z *Zip) SetDownloadBar(p interfaces.MultiProgressBar) {
	z.downloadBar = p.AddBar(int64(z.contentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"downloaded",
			),
		),
	)
}

func (z *Zip) IncrementDownloadBar(downloadedBytes int) {
	z.downloadBar.IncrBy(downloadedBytes)
}

func (z *Zip) SetExtractionBar(p interfaces.MultiProgressBar, zipSize int64) {
	z.extractionBar = p.AddBar(zipSize,
		mpb.BarQueueAfter(z.downloadBar),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"unzipped",
			),
		),
	)
}

func (z *Zip) IncrementExtractionBar(extractedBytes int) {
	z.extractionBar.IncrBy(extractedBytes)
}

func (z *Zip) Download(downloadedBytesChan chan<- int) (err error) {
	var (
		oopsBuilder = oops.Code("zip_download_failed").
			With("downloadedBytesChan", downloadedBytesChan)
	)

	if err = utils.Download(z.url, z.absPath, downloadedBytesChan); err != nil {
		err = oopsBuilder.
			With("url", z.url).
			With("absPath", z.absPath).
			Wrapf(err, "failed to download: %v", z.url)
		return
	}

	return
}

func (z *Zip) ExtractTo(extractToPath string, extractedBytesChan chan<- int) (err error) {
	var (
		oopsBuilder = oops.Code("zip_extract_to_failed").
			With("extractToPath", extractToPath).
			With("extractedBytesChan", extractedBytesChan)
	)

	if err = utils.Unzip(z.absPath, extractToPath, extractedBytesChan); err != nil {
		err = oopsBuilder.
			With("absPath", z.absPath).
			Wrapf(err, "failed to unzip: %v", z.absPath)
		return
	}

	return
}

func (z *Zip) Remove() (err error) {
	var (
		oopsBuilder = oops.Code("zip_remove_failed")
	)

	if err = os.Remove(z.absPath); err != nil {
		err = oopsBuilder.
			With("absPath", z.absPath).
			Wrapf(err, "failed to remove: %v", z.absPath)
		return
	}

	return
}
