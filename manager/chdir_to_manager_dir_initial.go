package manager

import (
	"os"

	"github.com/samber/oops"
)

func (manager *Manager) ChdirToManagerDirInitial() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDirWith")

	if err := os.Chdir(manager.Dir.Initial); err != nil {
		return oopsBuilder.
			With("initialDir", manager.Dir.Initial).
			Wrapf(err, "failed to change to initial dir")
	}

	return nil
}
