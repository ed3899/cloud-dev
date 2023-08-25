package manager

import (
	"os"

	"github.com/samber/oops"
)

// Changes the current working directory to the run dir
func (m *Manager) GoToDirRun() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("GoToDirRun")

	if err := os.Chdir(m.Path.Dir.Run); err != nil {
		return oopsBuilder.
			With("runDir", m.Path.Dir.Run).
			Wrapf(err, "failed to change to run dir")
	}

	return nil
}
