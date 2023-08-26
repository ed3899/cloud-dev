package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

// Returns a new environment for the given tool and cloud. pathToPackerManifest is only used for terraform,
// use "" for packer.
func NewEnvironment(
	tool iota.Tool,
	cloud iota.Cloud,
	pathToPackerManifest string,
) (*Environment[any], error) {
	oopsBuilder := oops.
		Code("NewEnvironment").
		With("tool", tool).
		With("cloud", cloud).
		With("pathToPackerManifest", pathToPackerManifest).
		In("manager").
		In("environment").
		Tags("Environment")

	switch tool {
	case iota.Packer:
		packerEnvironment, err := NewPackerEnvironment(cloud)
		if err != nil {
			return nil, oopsBuilder.
				With("cloud", cloud).
				Wrapf(err, "failed to create packer environment")
		}

		return &Environment[any]{
			Base:  packerEnvironment.Base,
			Cloud: packerEnvironment.Cloud,
		}, nil

	case iota.Terraform:
		terraformEnvironment, err := NewTerraformEnvironment(pathToPackerManifest, cloud)
		if err != nil {
			return nil, oopsBuilder.
				With("cloud", cloud).
				Wrapf(err, "failed to create terraform environment")
		}

		return &Environment[any]{
			Base:  terraformEnvironment.Base,
			Cloud: terraformEnvironment.Cloud,
		}, nil

	default:
		return nil, oopsBuilder.
			Errorf("unknown tool: %v", tool)
	}
}

type Environment[T any] struct {
	Base  T
	Cloud any
}
