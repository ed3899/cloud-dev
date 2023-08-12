package manager

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func UnsetCredentialsWith(
	osUnsetenv func(string) error,
) ForCloudGetter {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCredentialsWith")

	forManager := func(manager ICloudGetter[iota.Cloud]) error {
		managerCloudName := manager.Cloud().Name()

		switch manager.Cloud() {
		case iota.Aws:
			for key := range awsCredentials {
				if err := osUnsetenv(key); err != nil {
					return oopsBuilder.
						With("cloudName", managerCloudName).
						Wrapf(err, "failed to unset environment variable %s", key)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", managerCloudName).
				Errorf("unknown cloud: %#v", manager.Cloud())
		}

		return nil
	}

	return forManager
}
