package manager

import "github.com/samber/oops"

func ChangeToRunDirWith(
	osChdir func(string) error,
) ForDirGetter {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToRunDirWith")

	forManager := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Run()); err != nil {
			return oopsBuilder.
				With("runDir", manager.Dir().Run()).
				Wrapf(err, "failed to change to run dir")
		}

		return nil
	}

	return forManager
}
