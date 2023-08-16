package template

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewTemplate(
	tool iota.Tool,
	cloud iota.Cloud,
) (*TemplateFile, error) {
	oopsBuilder := oops.
		Code("NewTemplate").
		With("tool", tool).
		With("cloud", cloud)

	currentExecutablePath, err := os.Executable()
	if err != nil {
		return nil, oopsBuilder.
			Wrapf(err, "failed to get current executable path")
	}

	return &TemplateFile{
		Path: filepath.Join(
			filepath.Dir(currentExecutablePath),
			iota.Templates.Name(),
		),
	}, nil
}

type TemplateFile struct {
	Path string
}
