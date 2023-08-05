package environment

import (
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func PickPackerCloudEnvironment(cloud cloud.Cloud) (NewPackerCloudEnvironment NewPackerCloudEnvironmentF, err error) {
	var (
		oopsBuilder = oops.
			Code("PickPackerCloudEnvironment").
			With("cloud", cloud.Name)
	)

	switch cloud.Kind {
	case constants.Aws:
		NewPackerCloudEnvironment = NewPackerAwsEnvironment

	default:
		err = oopsBuilder.
			Errorf("Unsupported cloud kind: %s", cloud.Name)
		return
	}

	return
}

type PickPackerCloudEnvironmentF func(cloud.Cloud) (NewPackerCloudEnvironment NewPackerCloudEnvironmentF, err error)

type NewPackerCloudEnvironmentF func() PackerCloudEnvironmentI

type PackerCloudEnvironmentI interface {
	IsPackerCloudEnvironment() bool
}
