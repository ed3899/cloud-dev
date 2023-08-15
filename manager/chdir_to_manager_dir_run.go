package manager

import (
	"os"

	"github.com/samber/oops"
)

func ChDirToManagerDirRun(
	manager *Manager,
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ManagerDirRunChdirWith")

	if err := os.Chdir(manager.Dir.Run); err != nil {
		return oopsBuilder.
			With("runDir", manager.Dir.Run).
			Wrapf(err, "failed to change to run dir")
	}

	return nil
}
