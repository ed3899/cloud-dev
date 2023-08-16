package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) DeleteTemplate() error {
	oopsBuilder := oops.
		In("manager").
		Code("Delete")

	err := os.Remove(m.Path.Template.Merged)
	if err != nil {
		return oopsBuilder.Wrapf(err, "failed to delete merged template")
	}

	return nil
}
