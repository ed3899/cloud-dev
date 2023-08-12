package manager

import "github.com/samber/oops"

func ChangeToInitialDirWith(
	osChdir func(string) error,
) ForDirGetter {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	forManager := func(manager IDirGetter) error {
		if err := osChdir(manager.Dir().Initial()); err != nil {
			return oopsBuilder.
				With("initialDir", manager.Dir().Initial()).
				Wrapf(err, "failed to change to initial dir")
		}

		return nil
	}

	return forManager
}
