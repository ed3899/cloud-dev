package aws

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/tool"
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
		packerDirName      = tool.PACKER_NAME
		awsDirName         = cloud.AWS_NAME
		packerManifestName = tool.PACKER_MANIFEST_NAME

		absPath        string
		lastBuiltAmiId string
	)

	if absPath, err = filepath.Abs(filepath.Join(packerDirName, awsDirName, packerManifestName)); err != nil {
		err = oopsBuilder.
			With("packerDirName", packerDirName).
			With("awsDirName", awsDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", packerManifestName)
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
