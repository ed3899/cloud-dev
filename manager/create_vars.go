package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) CreateVars() (*os.File, error) {
	oopsBuilder := oops.
		Code("CreateVars")

	file, err := os.Create(m.Path.Vars)
	if err != nil {
		return nil, oopsBuilder.
			Wrapf(err, "failed to create vars file: %s", m.Path.Vars)
	}

	return file, nil
}
