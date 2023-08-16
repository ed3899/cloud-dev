package template

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func (t *TemplateFile) Create() error {
	oopsBuilder := oops.
		Code("Create")

	currentExecutablePath, err := os.Executable()
	if err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to get current executable path")
	}

	currentExecutableDir := filepath.Dir(currentExecutablePath)

	templatePath := func(templateName string) string {
		return filepath.Join(
			currentExecutableDir,
			iota.Templates.Name(),
			t.Tool.Name(),
			templateName,
		)
	}

	err = file.MergeFilesTo(
		t.Path,
		templatePath(t.Cloud.Template().Base),
		templatePath(t.Cloud.Template().Cloud),
	)
	if err != nil {
		return oopsBuilder.
			Wrapf(err, "failed to merge template files")
	}

	return nil
}
