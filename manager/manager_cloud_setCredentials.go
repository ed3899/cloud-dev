package manager

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ManagerCloudSetCredentialsWith(
	osSetenv func(string, string) error,
) ManagerCloudSetCredentials {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCloudCredentials")

	managerCloudSetCredentials := func(manager interfaces.IClone[Manager]) error {
		managerClone := manager.Clone()

		switch managerClone.Cloud() {
		case iota.Aws:
			for key, value := range awsCredentials {
				if err := osSetenv(key, value); err != nil {
					return oopsBuilder.
						With("cloudName", managerClone.Cloud().Name()).
						Wrapf(err, "failed to set environment variable %s to %s", key, value)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", managerClone.Cloud().Name()).
				Errorf("unknown cloud: %#v", managerClone.Cloud())
		}

		return nil
	}

	return managerCloudSetCredentials
}

type ManagerCloudSetCredentials func(interfaces.IClone[Manager]) error
