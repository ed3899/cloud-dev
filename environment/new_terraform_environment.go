package environment

import (
	"github.com/ed3899/kumo/common/iota"
	_manager "github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
)

func NewTerraformEnvironment(
	manager *_manager.Manager,
	cloud iota.Cloud,
) (*TerraformEnvironment, error) {
	oopsBuilder := oops.
		Code("NewTerraformEnvironment").
		With("cloud", cloud)

	general := NewTerraformGeneralEnvironment()

	switch cloud {
	case iota.Aws:
		aws, err := NewTerraformAwsEnvironment(manager)
		if err != nil {
			return nil, oopsBuilder.
				Wrapf(err, "failed to create aws environment")
		}

		return &TerraformEnvironment{
			General: general,
			Cloud:   aws,
		}, nil

	default:
		return nil, oopsBuilder.
			Errorf("unknown cloud: %v", cloud)
	}
}

type TerraformEnvironment struct {
	General *TerraformGeneralEnvironment
	Cloud   any
}
