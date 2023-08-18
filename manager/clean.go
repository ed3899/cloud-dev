package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func (m *Manager) CleanPacker() error {
	oopsBuilder := oops.
		Code("Clean").
		In("manager").
		Tags("Manager")

	if m.Tool.Iota() != iota.Packer {
		err := oopsBuilder.
			Errorf("tool is not packer")

		return err
	}

	err := os.Remove(m.Path.PackerManifest)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove packer manifest file")

		return err
	}

	err = os.Remove(m.Path.Executable)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove executable")

		return err
	}

	err = os.Remove(m.Path.Vars)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove vars")

		return err
	}

	err = os.RemoveAll(m.Path.Dir.Plugins)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to remove plugins directory")

		return err
	}

	return nil
}
