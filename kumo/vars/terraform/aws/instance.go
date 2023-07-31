package aws

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/samber/oops"
)

type HashicorpVars struct {
	file *os.File
}

func NewHashicorpVars() (hashicorpVars *HashicorpVars, err error) {
	var (
		oopsBuilder = oops.
				Code("new_hashicorp_vars_failed")

		varsFile          *os.File
		absPathToVarsFile string
	)

	if absPathToVarsFile, err = filepath.Abs(filepath.Join(dirs.PACKER_DIR_NAME, dirs.AWS_DIR_NAME, hashicorp_vars.PACKER_VARS_NAME)); err != nil {
		err = oopsBuilder.
			With("dirs.PACKER_DIR_NAME", dirs.PACKER_DIR_NAME).
			With("dirs.AWS_DIR_NAME", dirs.AWS_DIR_NAME).
			Wrapf(err, "Error occurred while crafting absolute path to %s", hashicorp_vars.PACKER_VARS_NAME)
		return
	}

	if varsFile, err = os.Create(absPathToVarsFile); err != nil {
		err = oopsBuilder.
			With("absPathToVarsFile", absPathToVarsFile).
			Wrapf(err, "Error occurred while creating %s", hashicorp_vars.PACKER_VARS_NAME)
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
