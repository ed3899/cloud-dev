package download

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	_manager "github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/url"
	"github.com/vbauerster/mpb/v8"
)

func NewDownload(
	currentExecutableDir string,
	urlContentLength int64,
	manager interfaces.IClone[*_manager.Manager],
) *Download {
	managerClone := manager.Clone()

	hashicorpUrl := url.BuildHashicorpUrl(
		managerClone.Tool,
		runtime.GOOS,
		runtime.GOARCH,
	)

	return &Download{
		Name: managerClone.Tool.Name(),
		Path: &Path{
			Zip: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				fmt.Sprintf("%s.zip", managerClone.Tool.Name()),
			),
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				managerClone.Tool.Name(),
				fmt.Sprintf("%s.exe", managerClone.Tool.Name()),
			),
		},
		Url:           hashicorpUrl,
		ContentLength: urlContentLength,
		Bar:           &Bar{},
	}
}

func (d *Download) Clone() *Download {
	return &Download{
		Name:          d.Name,
		Path:          d.Path.Clone(),
		Url:           d.Url,
		ContentLength: d.ContentLength,
		Bar:           d.Bar.Clone(),
	}
}

type Download struct {
	Name, Url     string
	ContentLength int64
	Path          *Path
	Bar           *Bar
}

func (p *Path) Clone() *Path {
	return &Path{
		Zip:        p.Zip,
		Executable: p.Executable,
	}
}

type Path struct {
	Zip        string
	Executable string
}

func (b *Bar) Clone() *Bar {
	return &Bar{
		Downloading: b.Downloading,
		Extracting:  b.Extracting,
	}
}

type Bar struct {
	Downloading, Extracting *mpb.Bar
}

type IIncrBy interface {
	IncrBy(int)
}

type IAddBar interface {
	AddBar(int64, ...any) *mpb.Bar
}
