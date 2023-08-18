package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func (m *Manager) Clean() error {
	oopsBuilder := oops.
		Code("Clean").
		In("manager").
		Tags("Manager")

	err := os.Remove(m.Path.Executable)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to remove executable")

		return err
	}

	err = os.Remove(m.Path.Vars)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to remove vars")

		return err
	}

	err = os.RemoveAll(m.Path.Dir.Plugins)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to remove plugins directory")

		return err
	}

	switch m.Tool.Iota() {
	case iota.Packer:
		err = os.Remove(m.Path.PackerManifest)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove packer manifest file")

			return err
		}

	case iota.Terraform:
		err = os.Remove(m.Path.Terraform.Lock)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove terraform lock file")

			return err
		}

		err = os.Remove(m.Path.Terraform.State)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove terraform state file")

			return err
		}

		err = os.Remove(m.Path.Terraform.Backup)
		if err != nil {
			err = oopsBuilder.
				Wrapf(err, "failed to remove terraform backup file")

			return err
		}

	default:
		err = oopsBuilder.
			Errorf("unknown tool: %#v", m.Tool)

		return err
	}

	return nil
}
