package template

import (
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

	pathToTemplate := func ()  {
		filepath.Join(
			iota.Template.Name()
		)
	}

	switch tool {
	case iota.Terraform:

		switch cloud {
		case iota.Aws:
			return &TemplateFile{
				Path: pathToTemplate,
			}, nil

		default:
			return nil, oopsBuilder.
				Errorf("cloud %s not supported", cloud)
		}

	case iota.Packer:

		switch cloud {
		case iota.Aws:

		}

	default:
		return nil, oopsBuilder.
			Errorf("tool %s not supported", tool)
	}
}

type TemplateFile struct {
	Path string
}
