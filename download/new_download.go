package download

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ed3899/kumo/common/iota"
	_manager "github.com/ed3899/kumo/manager"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func NewDownload(
	manager *_manager.Manager,
) (*Download, error) {
	oopsBuilder := oops.
		Code("NewDownload").
		In("download")

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current executable path")

		return nil, err
	}

	currentExecutableDir := filepath.Dir(currentExecutablePath)

	hashicorpUrl := url.BuildHashicorpUrl(
		manager.Tool.Name(),
		manager.Tool.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	contentLength, err := url.GetContentLength(hashicorpUrl)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get content length")

		return nil, err
	}

	return &Download{
		Name: manager.Tool.Name(),
		Path: &Path{
			Zip: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				fmt.Sprintf("%s.zip", manager.Tool.Name()),
			),
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				manager.Tool.Name(),
				fmt.Sprintf("%s.exe", manager.Tool.Name()),
			),
		},
		Url:           hashicorpUrl,
		ContentLength: contentLength,
		Bar:           &Bar{},
	}, nil
}

type Download struct {
	Name, Url     string
	ContentLength int64
	Path          *Path
	Bar           *Bar
}

type Path struct {
	Zip        string
	Executable string
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
