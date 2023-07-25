package binaries

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

func PackerBuildWorkflow() (err error) {
	var (
		cloud                    Cloud
		packer                   *Packer
		progress                 = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
		absPathToExecutableDir   string
		absPathToZipDir          string
		absPathToInitialLocation string
	)

	if packer, err = NewPacker(); err != nil {
		err = errors.Wrap(err, "Error occurred while creating new packer")
		return
	}

	if packer.IsInstalled() && packer.Zip.IsNotPresent() {
		return
	}

	// Start with a clean slate
	absPathToExecutableDir = filepath.Dir(packer.AbsPathToExecutable)
	absPathToZipDir = filepath.Dir(packer.Zip.AbsPath)

	if err := os.RemoveAll(absPathToExecutableDir); err != nil {
		err = errors.Wrapf(err, "Error occurred while removing %s", absPathToExecutableDir)
		return err
	}

	if err := os.RemoveAll(absPathToZipDir); err != nil {
		err = errors.Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return err
	}

	// Download
	if DownloadAndShowProgress(packer.Zip, progress); err != nil {
		err = errors.Wrapf(err, "Error occurred while downloading %s", packer.Zip.GetName())
		return
	}

	// Extract
	if ExtractAndShowProgress(packer.Zip, progress); err != nil {
		err = errors.Wrapf(err, "Error occurred while extracting %s", packer.Zip.GetName())
		return
	}

	// Set cloud
	if cloud, err = GetCloud(); err != nil {
		err = errors.Wrap(err, "Error occurred while getting cloud")
		return
	}

	// Set initial location
	if absPathToInitialLocation, err = os.Getwd(); err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return
	}

	// Change directory to packer run directory and defer changing back to initial location
	if err = os.Chdir(packer.AbsPathToRunDir); err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to packer run directory")
		return
	}
	defer func() {
		if err = os.Chdir(absPathToInitialLocation); err != nil {
			err = errors.Wrap(err, "Error occurred while changing directory back to initial location")
			return
		}
	}()

	// Initialize
	if err = packer.Init(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while initializing packer")
		return
	}

	// Build
	if err = packer.Build(cloud); err != nil {
		err = errors.Wrap(err, "Error occurred while building packer")
		return
	}

	return
}
