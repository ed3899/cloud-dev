package environment

import (
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/ed3899/kumo/constants"
)

type PickedCloudEnvironmentF func() CloudEnvironmentI

func PickPackerCloudEnvironment(cloud cloud.Cloud) (PickedCloudEnvironmentF PickedCloudEnvironmentF) {
	switch cloud.Kind {
	case constants.Aws:
		PickedCloudEnvironmentF = NewPackerAwsEnvironment
	default:
	}

	return
}
