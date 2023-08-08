package zip

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool"
	"github.com/ed3899/kumo/utils/url"

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

func New(t tool.Tool) (zip interfaces.ZipI, err error) {
	var (
		oopsBuilder = oops.Code("zip_new_failed")
		absPath     = filepath.Join(
			constants.DEPENDENCIES_DIR_NAME,
			t.Name,
			fmt.Sprintf(
				"%s.zip",
				t.Name,
			),
		)

		contentLength int64
	)

	if contentLength, err = url.GetContentLength(t.Url); err != nil {
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

func (z Zip) SetDownloadBar(p interfaces.ProgressBarAdder) {
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

func (z Zip) SetExtractionBar(p interfaces.ProgressBarAdder, zipSize int64) {
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
