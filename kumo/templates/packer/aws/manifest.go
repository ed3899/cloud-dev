package aws

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/templates"
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

	if absPath, err = filepath.Abs(filepath.Join(templates.PACKER_DIR_NAME, templates.AWS_DIR_NAME, templates.PACKER_MANIFEST_NAME)); err != nil {
		err = oopsBuilder.
			With("templates.PACKER_DIR_NAME", templates.PACKER_DIR_NAME).
			With("templates.AWS_DIR_NAME", templates.AWS_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", templates.PACKER_MANIFEST_NAME)
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
