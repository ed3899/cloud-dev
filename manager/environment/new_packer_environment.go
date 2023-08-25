package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewPackerEnvironment(
	cloud iota.Cloud,
) (*Environment[*PackerBaseEnvironment], error) {
	oopsBuilder := oops.
		Code("NewPackerEnvironment").
		In("manager").
		In("environment").
		With("cloud", cloud)

	base := NewPackerBaseEnvironment()

	switch cloud {
	case iota.Aws:
		aws := NewPackerAwsEnvironment()

		return &Environment[*PackerBaseEnvironment]{
			Cloud: aws,
			Base:  base,
		}, nil

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", cloud)

		return nil, err
	}
}
