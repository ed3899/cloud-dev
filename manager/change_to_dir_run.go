package manager

import "github.com/samber/oops"

func ChangeToDirRunWith(
	osChdir func(string) error,
) ChangeToDirRun {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToRunDirWith")

	changeToDirRun := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Run()); err != nil {
			return oopsBuilder.
				With("runDir", manager.Dir().Run()).
				Wrapf(err, "failed to change to run dir")
		}

		return nil
	}

	return changeToDirRun
}

type ChangeToDirRun func(IDirGetter) error
