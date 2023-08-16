package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewPackerEnvironment(
	cloud iota.Cloud,
) (*PackerEnvironment, error) {
	oopsBuilder := oops.
		Code("NewPackerEnvironment").
		With("cloud", cloud)

	general := NewPackerGeneralEnvironment()

	switch cloud {
	case iota.Aws:
		aws := NewPackerAwsEnvironment()

		return &PackerEnvironment{
			Cloud:   aws,
			General: general,
		}, nil

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", cloud.Name())

		return nil, err
	}
}

type PackerEnvironment struct {
	Cloud   any
	General *PackerGeneralEnvironment
}
