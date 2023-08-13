package manager

import "github.com/samber/oops"

func ChangeToDirRun(
	osChdir func(string) error,
	manager IDirGetter,
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToRunDirWith")

	if err := osChdir(manager.Dir().Run()); err != nil {
		return oopsBuilder.
			With("runDir", manager.Dir().Run()).
			Wrapf(err, "failed to change to run dir")
	}

	return nil
}
