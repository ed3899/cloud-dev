package manager

import (
	"os"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func UnsetManagerCloudCredentials(
	manager interfaces.IClone[*Manager],
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCloudCredentials")

	managerClone := manager.Clone()

	switch managerClone.Cloud {
	case iota.Aws:
		for key := range awsCredentials {
			if err := os.Unsetenv(key); err != nil {
				return oopsBuilder.
					With("cloudName", managerClone.Cloud.Name).
					Wrapf(err, "failed to unset environment variable %s", key)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", managerClone.Cloud.Name).
			Errorf("unknown cloud: %#v", managerClone.Cloud)
	}

	return nil

}
