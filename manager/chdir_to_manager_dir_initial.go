package manager

import (
	"os"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/samber/oops"
)

func ChdirToManagerDirInitial(manager interfaces.IClone[*Manager]) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	managerClone := manager.Clone()

	if err := os.Chdir(managerClone.Dir.Initial); err != nil {
		return oopsBuilder.
			With("initialDir", managerClone.Dir.Initial).
			Wrapf(err, "failed to change to initial dir")
	}

	return nil
}
