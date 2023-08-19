package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) DeletePluginsDir() error {
	oopsBuilder := oops.
		Code("DeletePluginsDir").
		In("manager").
		Tags("Manager")

	err := os.RemoveAll(m.Path.Dir.Plugins)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove plugins directory")

		return err
	}

	return nil
}
