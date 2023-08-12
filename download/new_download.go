package download

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func (d Download) Name() string {
	return d.name
}

type INameGetter interface {
	Name() string
}

func (d Download) Path() string {
	return d.path
}

type IPathGetter interface {
	Path() string
}

func (d Download) Url() string {
	return d.url
}

type IUrlGetter interface {
	Url() string
}

func (d Download) ContentLength() int64 {
	return d.contentLength
}

type IContentLengthGetter interface {
	ContentLength() int64
}

func (d Download) Bar() IBar {
	return d.bar
}

type IBarsGetter interface {
	Bar() IBar
}

type IDownload interface {
	INameGetter
	IPathGetter
	IUrlGetter
	IContentLengthGetter
	IBarsGetter
}

type Download struct {
	name, path, url string
	contentLength   int64
	bar             IBar
}

func (b Bars) Downloading() IIncrBy {
	return b.downloading
}

type IDownloading interface {
	Downloading() IIncrBy
}

func (b Bars) SetDownloading(mpbBar IIncrBy) {
	b.downloading = mpbBar
}

type ISetDownloading interface {
	SetDownloading(IIncrBy)
}

func (b Bars) Extracting() IIncrBy {
	return b.extracting
}

type IExtracting interface {
	Extracting() IIncrBy
}

func (b Bars) SetExtracting(mpbBar IIncrBy) {
	b.extracting = mpbBar
}

type ISetExtracting interface {
	SetExtracting(IIncrBy)
}

type IBar interface {
	IDownloading
	ISetDownloading
	ISetExtracting
	IExtracting
}

type Bars struct {
	downloading, extracting IIncrBy
}

type IIncrBy interface {
	IncrBy(int)
}

func SetDownloadBarWith(
	progress IAddBar,
) {
}

type IAddBar interface {
	AddBar(int64, ...any) IIncrBy
}

func (d Download) SetDownloadBar(p IAddBar) {
	d.bar.SetDownloading(
		p.AddBar(int64(d.contentLength),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"downloaded",
				),
			),
		),
	)
}

func (d Download) IncrementDownloadBar(downloadedBytes int) {
	d.bar.Downloading().IncrBy(downloadedBytes)
}

func (d Download) SetExtractionBar(p IAddBar, zipSize int64) {
	d.bar.SetExtracting(
		p.AddBar(zipSize,
			mpb.BarQueueAfter(d.bar.Downloading().(*mpb.Bar)),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.Name()),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"unzipped",
				),
			),
		),
	)

}

func (d *Download) IncrementExtractionBar(extractedBytes int) {
	d.bar.Extracting().IncrBy(extractedBytes)
}
