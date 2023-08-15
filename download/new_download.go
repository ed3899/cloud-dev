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
			path: Path{
				zip: filepath.Join(
					osExecutableDir,
					iota.Dependencies.Name(),
					fmt.Sprintf("%s.zip", manager.Tool().Name()),
				),
				executable: filepath.Join(
					osExecutableDir,
					iota.Dependencies.Name(),
					manager.Tool().Name(),
					fmt.Sprintf("%s.exe", manager.Tool().Name()),
				),
			},
			url: utilsUrlBuildHashicorpUrl(
				manager.Tool(),
				runtime.GOOS,
				runtime.GOARCH,
			),
			contentLength: contentLength,
			bar:           Bar{},
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

func (d Download) Path() Path {
	return d.path
}

type IPathSetter interface {
	SetPath(string) Download
}

func (d Download) SetPath(path IPath) Download {
	d.path = path.(Path)
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
		path:          d.path.Clone(),
		url:           d.url,
		contentLength: d.contentLength,
		bar:           d.bar.Clone(),
	}
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

type Download struct {
	name, url     string
	contentLength int64
	path          Path
	bar           Bar
}

type IZipGetter interface {
	Zip() string
}

func (p Path) Zip() string {
	return p.zip
}

type IZipSetter interface {
	SetZip(string) Path
}

func (p Path) SetZip(zip string) Path {
	p.zip = zip
	return p
}

type IExecutableGetter interface {
	Executable() string
}

func (p Path) Executable() string {
	return p.executable
}

type IExecutableSetter interface {
	SetExecutable(string) Path
}

func (p Path) SetExecutable(executable string) Path {
	p.executable = executable
	return p
}

func (p Path) Clone() Path {
	return Path{
		zip:        p.zip,
		executable: p.executable,
	}
}

type IPath interface {
	IZipGetter
	IZipSetter
	IExecutableGetter
	IExecutableSetter
	interfaces.IClone[Path]
}

type Path struct {
	zip        string
	executable string
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
