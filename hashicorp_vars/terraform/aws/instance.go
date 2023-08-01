package aws

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
)

type HashicorpVars struct {
	file *os.File
}

func NewHashicorpVars() (hashicorpVars *HashicorpVars, err error) {
	var (
		oopsBuilder = oops.
				Code("new_hashicorp_vars_failed")
		packerDirName = tool.PACKER_NAME
		awsDirName    = cloud.AWS_NAME
		varsFileName  = hashicorp_vars.PACKER_VARS_NAME

		varsFile          *os.File
		absPathToVarsFile string
	)

	if absPathToVarsFile, err = filepath.Abs(filepath.Join(packerDirName, awsDirName, varsFileName)); err != nil {
		err = oopsBuilder.
			With("packerDirName", packerDirName).
			With("awsDirName", awsDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", varsFileName)
		return
	}

	if varsFile, err = os.Create(absPathToVarsFile); err != nil {
		err = oopsBuilder.
			With("absPathToVarsFile", absPathToVarsFile).
			Wrapf(err, "Error occurred while creating %s", varsFileName)
		return
	}

	hashicorpVars = &HashicorpVars{
		file: varsFile,
	}

	return
}

func (hv *HashicorpVars) GetFile() (file *os.File) {
	return hv.file
}
