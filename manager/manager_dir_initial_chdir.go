package manager

import "github.com/samber/oops"

func ManagerDirInitialChdirWith(
	osChdir func(string) error,
) ManagerDirInitialChdir {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	managerDirInitialChdir := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Initial()); err != nil {
			return oopsBuilder.
				With("initialDir", manager.Dir().Initial()).
				Wrapf(err, "failed to change to initial dir")
		}

		return nil
	}

	return managerDirInitialChdir
}

type ManagerDirInitialChdir func(IDirGetter) error
