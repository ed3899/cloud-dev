package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) DeleteSshConfig() error {
	oopsBuilder := oops.
		Code("DeleteSshConfig")

	err := os.Remove(m.Path.SshConfig)
	if err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to delete ssh config file: %s", m.Path.SshConfig)
	}

	return nil
}
