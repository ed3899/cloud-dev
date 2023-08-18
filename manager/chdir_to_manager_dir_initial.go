package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) ChdirToManagerDirInitial() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ChangeToInitialDir")

	if err := os.Chdir(m.Path.Dir.Initial); err != nil {
		return oopsBuilder.
			With("initialDir", m.Path.Dir.Initial).
			Wrapf(err, "failed to change to initial dir")
	}

	return nil
}
