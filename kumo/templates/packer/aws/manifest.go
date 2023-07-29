package aws

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/packer_manifest"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Manifest struct {
	lastBuiltAmiId string
}

func (m *Manifest) GetLastBuiltAmiId() (lastBuiltAmiId string) {
	return m.lastBuiltAmiId
}

func NewManifest() (manifest *Manifest, err error) {
	var (
		oopsBuilder = oops.
				Code("new_manifest_failed")

		absPath        string
		lastBuiltAmiId string
	)

	if absPath, err = filepath.Abs(filepath.Join(dirs.PACKER_DIR_NAME, dirs.AWS_DIR_NAME, packer_manifest.NAME)); err != nil {
		err = oopsBuilder.
			With("dirs.PACKER_DIR_NAME", dirs.PACKER_DIR_NAME).
			With("dirs.AWS_DIR_NAME", dirs.AWS_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", packer_manifest.NAME)
		return
	}

	if lastBuiltAmiId, err = utils.GetLastBuiltAmiId(absPath); err != nil {
		err = oopsBuilder.
			With("absPath", absPath).
			Wrapf(err, "failed to get last built AMI ID")
		return
	}

	manifest = &Manifest{
		lastBuiltAmiId: lastBuiltAmiId,
	}

	return
}
