package structs

import (
	"path/filepath"

	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Manifest struct {
	lastBuiltAmiId string
}

func New() (manifest *Manifest, err error) {
	var (
		oopsBuilder = oops.
				Code("new_manifest_failed")
		packerDirName      = common_tool_constants.PACKER_NAME
		awsDirName         = common_cloud_constants.AWS_NAME
		packerManifestName = common_tool_constants.PACKER_MANIFEST_NAME

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

func (m *Manifest) LastBuiltAmiId() (lastBuiltAmiId string) {
	return m.lastBuiltAmiId
}
