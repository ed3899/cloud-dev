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

func NewDownload(
	options ...Option,
) (
	download *Download,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewDownload").
				With("options", options)

		opt Option
	)

	download = &Download{}
	for _, opt = range options {
		if err = opt(download); err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to apply option '%v'", opt)
			return
		}
	}

	return
}

func WithName(tool tool.Tool) Option {
	return func(download *Download) (err error) {
		download.Name = tool.Name

		return
	}
}

func WithAbsPath(tool tool.Tool) Option {
	return func(download *Download) (err error) {
		download.AbsPath = filepath.Join(
			constants.DEPENDENCIES_DIR_NAME,
			tool.Name,
			fmt.Sprintf(
				"%s.zip",
				tool.Name,
			),
		)

		return
	}
}

func WithUrl(tool tool.Tool) Option {
	return func(download *Download) (err error) {
		download.Url = tool.Url

		return
	}
}

func WithContentLength(tool tool.Tool, getContentLength url.GetContentLengthF) Option {
	var (
		oopsBuilder = oops.
				Code("WithContentLength").
				With("tool", tool)

		contentLength int64
	)

	return func(download *Download) (err error) {
		if contentLength, err = getContentLength(tool.Url); err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to get zip content length for '%s'", tool.Url)
			return
		}
		download.ContentLength = contentLength

		return
	}
}

type Option func(*Download) error

type Download struct {
	Name          string
	AbsPath       string
	Url           string
	ContentLength int64
	Progress      *Progress
}

type Progress struct {
	Downloading *mpb.Bar
	Extracting  *mpb.Bar
}

func (d *Download) SetDownloadBar(p interfaces.ProgressBarAdder) {
	d.Progress.Downloading = p.AddBar(int64(d.ContentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(d.Name),
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

func (d *Download) IncrementDownloadBar(downloadedBytes int) {
	d.Progress.Downloading.IncrBy(downloadedBytes)
}

func (d *Download) SetExtractionBar(p interfaces.ProgressBarAdder, zipSize int64) {
	d.Progress.Extracting = p.AddBar(zipSize,
		mpb.BarQueueAfter(d.Progress.Downloading),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(d.Name),
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

func (d *Download) IncrementExtractionBar(extractedBytes int) {
	d.Progress.Extracting.IncrBy(extractedBytes)
}
