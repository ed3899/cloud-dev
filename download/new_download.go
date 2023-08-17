package download

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func NewDownload(
	_manager *manager.Manager,
) (*Download, error) {
	oopsBuilder := oops.
		Code("NewDownload").
		In("download").
		Tags("Download").
		With("manager", _manager)

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current executable path")

		return nil, err
	}

	currentExecutableDir := filepath.Dir(currentExecutablePath)

	hashicorpUrl := url.BuildHashicorpUrl(
		_manager.Tool.Name(),
		_manager.Tool.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	contentLength, err := url.GetContentLength(hashicorpUrl)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get content length")

		return nil, err
	}

	progress := mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithAutoRefresh(), mpb.WithWidth(64))

	return &Download{
		Name: _manager.Tool.Name(),
		Path: &Path{
			Zip: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				fmt.Sprintf("%s.zip", _manager.Tool.Name()),
			),
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				_manager.Tool.Name(),
				fmt.Sprintf("%s.exe", _manager.Tool.Name()),
			),
		},
		Url:           hashicorpUrl,
		ContentLength: contentLength,
		Bar:           &Bar{},
		Progress:      progress,
	}, nil
}

type Download struct {
	Name, Url     string
	ContentLength int64
	Path          *Path
	Progress      *mpb.Progress
	Bar           *Bar
}

type Path struct {
	Zip        string
	Executable string
}

type Bar struct {
	Downloading, Extracting *mpb.Bar
}
