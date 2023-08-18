package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) DeletePackerManifest() error {
	oopsBuilder := oops.
		Code("DeletePackerManifest").
		In("manager").
		Tags("Manager")

	err := os.Remove(m.Path.PackerManifest)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove packer manifest %s", m.Path.PackerManifest)

		return err
	}

	return nil
}
