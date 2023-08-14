package manager

import "github.com/samber/oops"

func ChangeToDirInitialWith(
	osChdir func(string) error,
) ChangeToDirInitial {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	changeToDirInitial := func(manager IDirGetter[IInitialGetter]) error {
		if err := osChdir(manager.Dir().Initial()); err != nil {
			return oopsBuilder.
				With("initialDir", manager.Dir().Initial()).
				Wrapf(err, "failed to change to initial dir")
		}

		return nil
	}

	return changeToDirInitial
}

type ChangeToDirInitial func(IDirGetter[IInitialGetter]) error
