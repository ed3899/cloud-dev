package manager

import (
	"os"

	"github.com/samber/oops"
)

// Deletes the vars file.
func (m *Manager) DeleteVars() error {
	oopsBuilder := oops.
		Code("DeleteVars").
		In("manager").
		Tags("Manager")

	err := os.RemoveAll(m.Path.Vars)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove vars directory")

		return err
	}

	return nil
}
