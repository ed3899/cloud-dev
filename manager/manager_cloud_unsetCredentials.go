package manager

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ManagerCloudUnsetCredentialsWith(
	osUnsetenv func(string) error,
) ManagerCloudUnsetCredentials {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCloudCredentials")

	managerCloudUnsetCredentials := func(manager interfaces.IClone[Manager]) error {
		managerClone := manager.Clone()

		switch managerClone.Cloud() {
		case iota.Aws:
			for key := range awsCredentials {
				if err := osUnsetenv(key); err != nil {
					return oopsBuilder.
						With("cloudName", managerClone.Cloud().Name()).
						Wrapf(err, "failed to unset environment variable %s", key)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", managerClone.Cloud().Name()).
				Errorf("unknown cloud: %#v", managerClone.Cloud())
		}

		return nil
	}

	return managerCloudUnsetCredentials
}

type ManagerCloudUnsetCredentials func(interfaces.IClone[Manager]) error
