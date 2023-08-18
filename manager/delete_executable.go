package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) DeleteExecutable() error {
	oopsBuilder := oops.
		Code("DeleteExecutable").
		In("manager").
		Tags("Manager")

	err := os.Remove(m.Path.Executable)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove executable %s", m.Path.Executable)

		return err
	}

	return nil
}
