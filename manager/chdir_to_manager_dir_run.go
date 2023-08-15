package manager

import (
	"os"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/samber/oops"
)

func ChDirToManagerDirRun(
	manager interfaces.IClone[*Manager],
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ManagerDirRunChdirWith")

	managerClone := manager.Clone()

	if err := os.Chdir(managerClone.Dir.Run); err != nil {
		return oopsBuilder.
			With("runDir", managerClone.Dir.Run).
			Wrapf(err, "failed to change to run dir")
	}

	return nil
}
