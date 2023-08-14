package manager

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/samber/oops"
)

func ManagerDirRunChdirWith(
	osChdir func(string) error,
) ManagerDirRunChdir[Manager] {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ManagerDirRunChdirWith")

	managerDirRunChdir := func(manager interfaces.IClone[Manager]) error {
		managerClone := manager.Clone()

		if err := osChdir(managerClone.Dir().Run()); err != nil {
			return oopsBuilder.
				With("runDir", managerClone.Dir().Run()).
				Wrapf(err, "failed to change to run dir")
		}

		return nil
	}

	return managerDirRunChdir
}

type ManagerDirRunChdir[M IManager] func(interfaces.IClone[Manager]) error
