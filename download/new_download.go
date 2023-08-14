package download

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	_manager "github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

type ITool interface {
	Name() string
	Version() string
}

func NewDownloadWith(
	osExecutable func() (string, error),
	utilsUrlBuildHashicorpUrl func(ITool, string, string) string,
	utilsUrlGetContentLength func(string) (int64, error),
) NewDownload {
	oopsBuilder := oops.
		In("download").
		Code("NewDownloadWith")

	newDownload := func(manager _manager.Manager) (Download, error) {
		osExecutablePath, err := osExecutable()
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to get executable path")
			return Download{}, err
		}

		osExecutableDir := filepath.Dir(osExecutablePath)

		url := utilsUrlBuildHashicorpUrl(
			manager.Tool(),
			runtime.GOOS,
			runtime.GOARCH,
		)

		contentLength, err := utilsUrlGetContentLength(url)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to get content length")
			return Download{}, err
		}

		return Download{
			name: manager.Tool().Name(),
			path: filepath.Join(
				osExecutableDir,
				iota.Dependencies.Name(),
				fmt.Sprintf("%s.zip", manager.Tool().Name()),
			),
			url: utilsUrlBuildHashicorpUrl(
				manager.Tool(),
				runtime.GOOS,
				runtime.GOARCH,
			),
			contentLength: contentLength,
		}, nil
	}

	return newDownload
}

type NewDownload func(_manager.Manager) (Download, error)

type INameGetter interface {
	Name() string
}

func (d Download) Name() string {
	return d.name
}

type INameSetter interface {
	SetName(string) Download
}

func (d Download) SetName(name string) Download {
	d.name = name
	return d
}

type IPathGetter interface {
	Path() string
}

func (d Download) Path() string {
	return d.path
}

type IPathSetter interface {
	SetPath(string) Download
}

func (d Download) SetPath(path string) Download {
	d.path = path
	return d
}

type IUrlGetter interface {
	Url() string
}

func (d Download) Url() string {
	return d.url
}

type IUrlSetter interface {
	SetUrl(string) Download
}

func (d Download) SetUrl(url string) Download {
	d.url = url
	return d
}

type IContentLengthGetter interface {
	ContentLength() int64
}

func (d Download) ContentLength() int64 {
	return d.contentLength
}

type IContentLengthSetter interface {
	SetContentLength(int64) Download
}

func (d Download) SetContentLength(contentLength int64) Download {
	d.contentLength = contentLength
	return d
}

type IBarGetter interface {
	Bar() Bar
}

func (d Download) Bar() Bar {
	return d.bar
}

type IBarSetter interface {
	SetBar(IBar) Download
}

func (d Download) SetBar(bar IBar) Download {
	d.bar = bar.(Bar)
	return d
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

type Download struct {
	name, path, url string
	contentLength   int64
	bar             Bar
}

type IDownload interface {
	INameGetter
	INameSetter
	IPathGetter
	IPathSetter
	IUrlGetter
	IUrlSetter
	IContentLengthGetter
	IContentLengthSetter
	IBarGetter
	IBarSetter
	interfaces.IClone[Download]
}

type IDownloadingGetter interface {
	Downloading() *mpb.Bar
}

func (b Bar) Downloading() *mpb.Bar {
	return b.downloading
}

type IDownloadingSetter interface {
	SetDownloading(IIncrBy) Bar
}

func (b Bar) SetDownloading(mpbBar IIncrBy) Bar {
	b.downloading = mpbBar.(*mpb.Bar)
	return b
}

type IExtractingGetter interface {
	Extracting() IIncrBy
}

func (b Bar) Extracting() IIncrBy {
	return b.extracting
}

type IExtractingSetter interface {
	SetExtracting(IIncrBy) Bar
}

func (b Bar) SetExtracting(mpbBar IIncrBy) Bar {
	b.extracting = mpbBar.(*mpb.Bar)
	return b
}

func (b Bar) Clone() Bar {
	return Bar{
		downloading: b.downloading,
		extracting:  b.extracting,
	}
}

type Bar struct {
	downloading, extracting *mpb.Bar
}

type IBar interface {
	IDownloadingGetter
	IDownloadingSetter
	IExtractingGetter
	IExtractingSetter
	interfaces.IClone[Bar]
}

type IIncrBy interface {
	IncrBy(int)
}

type IAddBar interface {
	AddBar(int64, ...any) *mpb.Bar
}
