package manager

import (
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func (m *Manager) CreateTemplate() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("CreateTemplate")

	err := file.MergeFilesTo(m.Path.Template.Merged, m.Path.Template.Cloud, m.Path.Template.Base)
	if err != nil {
		return oopsBuilder.Wrapf(err, "failed to merge files")
	}

	return nil
}
