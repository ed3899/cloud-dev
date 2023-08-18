package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) ChDirToManagerDirRun() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ManagerDirRunChdir")

	if err := os.Chdir(m.Path.Dir.Run); err != nil {
		return oopsBuilder.
			With("runDir", m.Path.Dir.Run).
			Wrapf(err, "failed to change to run dir")
	}

	return nil
}
