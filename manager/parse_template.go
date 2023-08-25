package manager

import (
	"text/template"

	"github.com/samber/oops"
)

// Parses the merged template.
func (m *Manager) ParseTemplate() (*template.Template, error) {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("ParseTemplate")

	template, err := template.ParseFiles(m.Path.Template.Merged)
	if err != nil {
		return nil, oopsBuilder.Wrapf(err, "failed to parse template")
	}

	return template, nil
}
