package binaries

import (
	"sync"

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

}

func (t *Terraform2) Down() (err error) {

}

type Packer2I interface {
	Build() (err error)
}

type Packer2 struct {
	ID   Tool
	Name string
	Path string
	Zip  *Zip
}

func NewPacker() (packer *Packer2, err error) {
	name := "packer"
	version := utils.GetLatestPackerVersion()
	os, arch := utils.GetCurrentHostSpecs()
	depDirName := utils.GetDependenciesDirName()
	binaryPath, err := utils.CreateBinaryPath([]string{depDirName, name, "packer.exe"})
	zipPath, err := utils.CreateZipPath([]string{depDirName, name, "packer.zip"})
	url := utils.CreateHashicorpURL(name, version, os, arch)
	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %v", url)
		return
	}
	bwg := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&bwg), mpb.WithWidth(60), mpb.WithAutoRefresh())
	downloadBar := progress.AddBar(int64(contentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"downloaded",
			),
		),
	)

	packer = &Packer2{
		ID:   PackerID,
		Name: name,
		Path: binaryPath,
		Zip: &Zip{
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
			DownloadBar:   downloadBar,
		},
	}
}

func (p *Packer2) Build() (err error) {
	p.Zip.Extract(p.Path)
}

type ZipI interface {
	Download() (err error)
	Extract() (err error)
	Remove() (err error)
}

type Zip struct {
	Path          string   // will be crafted by util
	URL           string   // will be crafted by util
	ContentLength int64    // will be crafted by util
	DownloadBar   *mpb.Bar // will be crafted by util
	ExtractionBar *mpb.Bar // wil be crafted by util
}

func (z *Zip) Download() (err error) {
}

func (z *Zip) Extract(path string) (err error) {
}

func (z *Zip) Remove() (err error) {
}
