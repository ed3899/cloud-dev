package binaries

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

type Terraform2I interface {
	Up() (err error)
	Down() (err error)
}

type Terraform2 struct {
	Path string
	Zip  *Zip
}

func (t *Terraform2) Up() (err error) {
	return
}

func (t *Terraform2) Down() (err error) {
	return
}

type Packer2I interface {
	Build() (err error)
}

type Packer2 struct {
	ID   Tool
	Name string
	AbsPathToExecutable string
	Zip  *Zip
}

func NewPacker() (packer *Packer2, err error) {
	const (
		name = "packer"
		version = "1.9.1"
	)

	var (
		executableName = fmt.Sprintf("%s.exe", name)
		zipName = fmt.Sprintf("%s.zip", name)
		os, arch = utils.GetCurrentHostSpecs()
		url = utils.CreateHashicorpURL(name, version, os, arch)
		depDirName = utils.GetDependenciesDirName()
	)

	absPathToExecutable, err := filepath.Abs(filepath.Join(depDirName, name, executableName))
	if err != nil {
		err = errors.Wrapf(err, "failed to create binary path to: %v", executableName)
		return
	}
	zipPath, err := filepath.Abs(filepath.Join(depDirName, name, zipName))
	if err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %v", zipName)
		return
	}
	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %v", url)
		return
	}

	packer = &Packer2{
		ID:   PackerID,
		AbsPathToExecutable: absPathToExecutable,
		Zip: &Zip{
			Name:          zipName,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}

	return
}

func (p *Packer2) Build() (err error) {
	return
}

type ZipI interface {
	Download() (err error)
	Extract() (err error)
	Remove() (err error)
}

type Zip struct {
	Name          string
	Path          string
	URL           string
	ContentLength int64
	DownloadBar   *mpb.Bar
	ExtractionBar *mpb.Bar
}

func (z *Zip) SetDownloadBar(p *mpb.Progress) {
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

func (z *Zip) SetExtractionBar(p *mpb.Progress) (err error) {
	if utils.FileNotPresent(z.Path) {
		err = errors.New("zip file not present")
		return
	}

	zipSize, err := utils.GetZipSize(z.Path)
	if err != nil {
		err = errors.Wrapf(err, "failed to get zip size for: %v", z.Path)
		return
	}

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

	return
}

func (z *Zip) Download(downloadedBytesChan chan<- int) (err error) {
	if utils.FilePresent(z.Path) {
		err = errors.New("zip file already present")
		return
	}

	if z.DownloadBar == nil {
		err = errors.New("download bar not set")
		return
	}

	if err = utils.Download(z.URL, z.Path, downloadedBytesChan); err != nil {
		err = errors.Wrapf(err, "failed to download: %v", z.URL)
		return
	}
	return
}

func (z *Zip) Extract(extractToPath string, extractedBytesChan chan<- int) (err error) {
	if utils.FileNotPresent(z.Path) {
		err = errors.New("zip file not present")
		return
	}

	if z.ExtractionBar == nil {
		err = errors.New("extraction bar not set")
		return
	}

	if err = utils.Unzip(z.Path, extractToPath, extractedBytesChan); err != nil {
		err = errors.Wrapf(err, "failed to unzip: %v", z.Path)
		return
	}
	return
}

func (z *Zip) Remove() (err error) {
	if err = os.Remove(z.Path); err != nil {
		err = errors.Wrapf(err, "failed to remove: %v", z.Path)
		return
	}
	return
}
