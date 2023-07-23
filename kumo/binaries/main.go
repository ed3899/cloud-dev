package binaries

import (
	"fmt"
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

	packer = &Packer2{
		ID:   PackerID,
		Name: name,
		Path: binaryPath,
		Zip: &Zip{
			Name : name,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}

	return
}


func (p *Packer2) Build() (err error) {
}

type ZipI interface {
	Download() (err error)
	Extract() (err error)
	Remove() (err error)
}

type Zip struct {
	Name string
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

	barName := fmt.Sprintf("%s.zip", z.Name)

	z.ExtractionBar = p.AddBar(zipSize,
		mpb.BarQueueAfter(z.DownloadBar),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(barName),
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

func (z *Zip) Download() (err error) {
	
}

func (z *Zip) Extract(path string) (err error) {
}

func (z *Zip) Remove() (err error) {
}
