package aws

import (
	"os"

	"github.com/samber/oops"
)

type Vars struct {
	file *os.File
}

func NewVars() (vars *Vars, err error) {
	var (
		oopsBuilder = oops.
			Code("new_vars_failed")

		varsFile *os.File
	)

	if varsFile, err = os.Create(hv.AbsPath); err != nil {
		return errors.Wrapf(err, "Error occurred while creating %s", hv.AbsPath)
	}

	return
}

func (v *Vars) GetFile() (file *os.File) {
	return v.file
}
