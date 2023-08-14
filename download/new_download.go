package download

import (
	"github.com/ed3899/kumo/common/interfaces"
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

func (d Download) Bar() Bar {
	return d.bar
}

type IBarGetter[Bar any] interface {
	Bar() Bar
}

func (d Download) Clone() Download {
	return Download{
		name:          d.name,
		path:          d.path,
		url:           d.url,
		contentLength: d.contentLength,
		bar:           d.bar.Clone(),
	}
}

type IDownload interface {
	INameGetter
	IPathGetter
	IUrlGetter
	IContentLengthGetter
	IBarGetter[Bar]
	interfaces.IClone[Download]
}

type Download struct {
	name, path, url string
	contentLength   int64
	bar             Bar
}

func (b Bar) Downloading() IIncrBy {
	return b.downloading
}

type IDownloadingGetter interface {
	Downloading() IIncrBy
}

func (b Bar) SetDownloading(mpbBar IIncrBy) {
	b.downloading = mpbBar
}

type IDownloadingSetter interface {
	SetDownloading(IIncrBy)
}

func (b Bar) Extracting() IIncrBy {
	return b.extracting
}

type IExtractingGetter interface {
	Extracting() IIncrBy
}

func (b Bar) SetExtracting(mpbBar IIncrBy) {
	b.extracting = mpbBar
}

type IExtractingSetter interface {
	SetExtracting(IIncrBy)
}

func (b Bar) Clone() Bar {
	return Bar{
		downloading: b.downloading,
		extracting:  b.extracting,
	}
}

type IBar interface {
	IDownloadingGetter
	IDownloadingSetter
	IExtractingGetter
	IExtractingSetter
	interfaces.IClone[Bar]
}

type Bar struct {
	downloading, extracting IIncrBy
}

type IIncrBy interface {
	IncrBy(int)
}

type IAddBar interface {
	AddBar(int64, ...any) IIncrBy
}

func (d *Download) IncrementExtractionBar(extractedBytes int) {
	d.bar.Extracting().IncrBy(extractedBytes)
}
