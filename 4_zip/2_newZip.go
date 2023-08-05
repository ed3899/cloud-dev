package zip

import (
	"fmt"
	"os"
	"path/filepath"

	constants "github.com/ed3899/kumo/0_constants"
	utils "github.com/ed3899/kumo/1_utils"
	tool "github.com/ed3899/kumo/3_tool"

	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Zip struct {
	Name          string
	AbsPath       string
	Url           string
	ContentLength int64
	DownloadBar   *mpb.Bar
	ExtractionBar *mpb.Bar
}

func New(t tool.Tool) (zip ZipI, err error) {
	var (
		oopsBuilder = oops.Code("zip_new_failed")
		absPath     = filepath.Join(
			constants.DEPENDENCIES_DIR_NAME,
			t.Name,
			fmt.Sprintf("%s.zip",
				t.Name),
		)

		contentLength int64
	)

	if contentLength, err = utils.GetContentLength(t.Url); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get zip content length")
		return
	}

	zip = Zip{
		Name:          filepath.Base(absPath),
		AbsPath:       absPath,
		Url:           t.Url,
		ContentLength: contentLength,
	}

	return
}

func (z Zip) SetDownloadBar(p MultiProgressBar) {
	z.DownloadBar = p.AddBar(int64(z.ContentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.Name),
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

func (z Zip) IncrementDownloadBar(downloadedBytes int) {
	z.DownloadBar.IncrBy(downloadedBytes)
}

func (z Zip) SetExtractionBar(p MultiProgressBar, zipSize int64) {
	z.ExtractionBar = p.AddBar(zipSize,
		mpb.BarQueueAfter(z.DownloadBar),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.Name),
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

func (z Zip) IncrementExtractionBar(extractedBytes int) {
	z.ExtractionBar.IncrBy(extractedBytes)
}

func (z Zip) Download(downloadedBytesChan chan<- int) (err error) {
	var (
		oopsBuilder = oops.Code("zip_download_failed").
			With("downloadedBytesChan", downloadedBytesChan)
	)

	if err = utils.Download(z.Url, z.AbsPath, downloadedBytesChan); err != nil {
		err = oopsBuilder.
			With("url", z.Url).
			With("absPath", z.AbsPath).
			Wrapf(err, "failed to download: %v", z.Url)
		return
	}

	return
}

func (z Zip) ExtractTo(extractToPath string, extractedBytesChan chan<- int) (err error) {
	var (
		oopsBuilder = oops.Code("zip_extract_to_failed").
			With("extractToPath", extractToPath).
			With("extractedBytesChan", extractedBytesChan)
	)

	if err = utils.Unzip(z.AbsPath, extractToPath, extractedBytesChan); err != nil {
		err = oopsBuilder.
			With("absPath", z.AbsPath).
			Wrapf(err, "failed to unzip: %v", z.AbsPath)
		return
	}

	return
}

func (z Zip) Remove() (err error) {
	var (
		oopsBuilder = oops.Code("zip_remove_failed")
	)

	if err = os.Remove(z.AbsPath); err != nil {
		err = oopsBuilder.
			With("absPath", z.AbsPath).
			Wrapf(err, "failed to remove: %v", z.AbsPath)
		return
	}

	return
}
