package workflows

import (
	"os"

	"github.com/samber/oops"
)

type PackerVars struct {
	File    *os.File
}

func NewPackerVars(template *PackerMergedTemplate) (packerVars *PackerVars, err error) {
	var (
		oopsBuilder = oops.
			Code("new_packer_vars_failed").
			With("template.Name", template.Name).
			With("template.AbsPath", template.AbsPath).
			With("template.Environment", template.Environment)
	)

	return
}
