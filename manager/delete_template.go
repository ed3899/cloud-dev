package manager

import (
	"os"

	"github.com/samber/oops"
)

// Deletes the merged template file.
func (m *Manager) DeleteTemplate() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("DeleteTemplate")

	err := os.Remove(m.Path.Template.Merged)
	if err != nil {
		return oopsBuilder.Wrapf(err, "failed to delete merged template")
	}

	return nil
}
