package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) Clean() error {
	oopsBuilder := oops.
		Code("Clean").
		In("manager").
		Tags("Manager")

	err := os.Remove(m.Path.PackerManifest)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove packer manifest file")

		return err
	}

	err = os.Remove(m.Path.Executable)

	return nil
}
