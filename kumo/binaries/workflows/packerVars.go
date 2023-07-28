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
			With("template.GetName()", template.GetName()).
			With("template.GetEnvironment()", template.GetEnvironment())

		packerVarsFile *os.File
		absPathToPackerVarsFile string
	)



	return
}
