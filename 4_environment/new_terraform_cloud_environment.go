package environment

import (
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

type NewTerraformCloudEnvironmentF func(pickedAmiId string) CloudEnvironmentI

func PickTerraformCloudEnvironment(cloud cloud.Cloud) (NewTerraformCloudEnvironment NewTerraformCloudEnvironmentF, err error) {
	var (
		oopsBuilder = oops.
			Code("NewTerraformCloudEnvironment").
			With("cloud", cloud.Name)
	)

	switch cloud.Kind {
	case constants.Aws:
		NewTerraformCloudEnvironment = NewTerraformAwsEnvironment
	default:
		err = oopsBuilder.
			Errorf("cloud not supported")
	}

	return

}

type TerraformCloudEnvironmentI interface {
	IsTerraformCloudEnvironment() bool
}
