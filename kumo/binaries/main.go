package binaries

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
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
	depDirName := utils.GetDependenciesDirName()
	binaryPath, err := utils.CreateBinaryPath([]string{depDirName, "packer", "packer.exe"})


	zipPath, err := filepath.Abs(filepath.Join("deps", "packer.zip"))
	if err != nil {
		err = errors.Wrap(err, "failed to get absolute path for packer zip")
		return
	}

	url := utils.CreateHashicorpURL()

	packer = &Packer2{
		ID: PackerID,
		Name: "packer",
		Path: binaryPath,
		Zip: &Zip{

		}
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
	Path          string // will be crafted by util
	URL           string // will be crafted by util
	ContentLength int64  // will be crafted by util
	DownloadBar   *mpb.Bar // will be crafted by util
	ExtractionBar *mpb.Bar  // wil be crafted by util
}

func (z *Zip) Download() (err error) {
}

func (z *Zip) Extract(path string) (err error) {
}

func (z *Zip) Remove() (err error) {
}
