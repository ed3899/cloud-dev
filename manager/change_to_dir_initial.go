package manager

import "github.com/samber/oops"

func ChangeToDirInitial(
	osChdir func(string) error,
	manager IDirGetter,
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	if err := osChdir(manager.Dir().Initial()); err != nil {
		return oopsBuilder.
			With("initialDir", manager.Dir().Initial()).
			Wrapf(err, "failed to change to initial dir")
	}

	return nil
}
