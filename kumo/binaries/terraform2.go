package binaries

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type Terraform2I interface {
	Init(Cloud) error
	Up(Cloud) error
	Destroy(Cloud) error
}

type Terraform2 struct {
	ID                  Tool
	AbsPathToExecutable string
	AbsPathToRunDir     string
	Zip                 *Zip
}

func NewTerraform() (terraform *Terraform2, err error) {
	const (
		name    = "terraform"
		version = "0.14.7"
	)

	var (
		executableName      = fmt.Sprintf("%s.exe", name)
		zipName             = fmt.Sprintf("%s.zip", name)
		os, arch            = utils.GetCurrentHostSpecs()
		url                 = utils.CreateHashicorpURL(name, version, os, arch)
		depDirName          = utils.GetDependenciesDirName()
		absPathToExecutable string
		absPathToRunDir     string
		contentLength       int64
		zipPath             string
	)

	if absPathToExecutable, err = filepath.Abs(filepath.Join(depDirName, name, executableName)); err != nil {
		err = errors.Wrapf(err, "failed to create executable path to: %s", executableName)
		return
	}

	if absPathToRunDir, err = filepath.Abs(name); err != nil {
		err = errors.Wrapf(err, "failed to create run path to: %s", name)
		return
	}

	if zipPath, err = filepath.Abs(filepath.Join(depDirName, name, zipName)); err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %s", zipName)
		return
	}

	if contentLength, err = utils.GetContentLength(url); err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	terraform = &Terraform2{
		ID:                  TerraformID,
		AbsPathToExecutable: absPathToExecutable,
		AbsPathToRunDir:     absPathToRunDir,
		Zip: &Zip{
			Name:          zipName,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}
	return
}

func (t *Terraform2) Init(cloud Cloud) (err error) {
	return
}

func (t *Terraform2) Up(cloud Cloud) (err error) {
	return
}

func (t *Terraform2) Destroy(cloud Cloud) (err error) {
	return
}
