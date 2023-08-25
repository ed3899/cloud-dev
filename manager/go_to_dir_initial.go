package manager

import (
	"os"

	"github.com/samber/oops"
)

// Changes the current working directory to the initial dir
func (m *Manager) GoToDirInitial() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("GoToDirInitial")

	if err := os.Chdir(m.Path.Dir.Initial); err != nil {
		return oopsBuilder.
			With("initialDir", m.Path.Dir.Initial).
			Wrapf(err, "failed to change to initial dir")
	}

	return nil
}
