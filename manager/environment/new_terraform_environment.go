package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

// Creates a new Terraform environment for the given cloud.
func NewTerraformEnvironment(
	pathToPackerManifest string,
	cloud iota.Cloud,
) (*Environment[*TerraformBaseEnvironment], error) {
	oopsBuilder := oops.
		Code("NewTerraformEnvironment").
		With("cloud", cloud).
		With("pathToPackerManifest", pathToPackerManifest).
		In("manager").
		In("environment")

	base := NewTerraformBaseEnvironment()

	switch cloud {
	case iota.Aws:
		aws, err := NewTerraformAwsEnvironment(pathToPackerManifest)
		if err != nil {
			return nil, oopsBuilder.
				Wrapf(err, "failed to create aws environment")
		}

		return &Environment[*TerraformBaseEnvironment]{
			Base:  base,
			Cloud: aws,
		}, nil

	default:
		return nil, oopsBuilder.
			Errorf("unknown cloud: %v", cloud)
	}
}
