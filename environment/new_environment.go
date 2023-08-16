package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewEnvironment(
	tool iota.Tool,
	cloud iota.Cloud,
) error {
	oopsBuilder := oops.
		Code("NewEnvironment").
		With("tool", tool).
		With("cloud", cloud)

	switch tool {
	case iota.Terraform:
	case iota.Packer:
		general := NewPackerGeneralEnvironment()
	default:
		err := oopsBuilder.
			Errorf("unknown tool: %s", tool)

		return err
	}

	return nil
}
