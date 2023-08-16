package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewTerraformEnvironment(
	pathToPackerManifest string,
	cloud iota.Cloud,
) (*Environment[*TerraformGeneralEnvironment], error) {
	oopsBuilder := oops.
		Code("NewTerraformEnvironment").
		With("cloud", cloud)

	general := NewTerraformGeneralEnvironment()

	switch cloud {
	case iota.Aws:
		aws, err := NewTerraformAwsEnvironment(pathToPackerManifest, cloud)
		if err != nil {
			return nil, oopsBuilder.
				Wrapf(err, "failed to create aws environment")
		}

		return &Environment[*TerraformGeneralEnvironment]{
			General: general,
			Cloud:   aws,
		}, nil

	default:
		return nil, oopsBuilder.
			Errorf("unknown cloud: %v", cloud)
	}
}
