package manager

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func UnsetCloudCredentialsWith(
	osUnsetenv func(string) error,
) UnsetCloudCredentials {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCloudCredentials")

	unsetCloudCredentials := func( manager ICloudGetter[iota.Cloud]) error {
		switch manager.Cloud() {
		case iota.Aws:
			for key := range awsCredentials {
				if err := osUnsetenv(key); err != nil {
					return oopsBuilder.
						With("cloudName", manager.Cloud().Name()).
						Wrapf(err, "failed to unset environment variable %s", key)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", manager.Cloud().Name()).
				Errorf("unknown cloud: %#v", manager.Cloud())
		}

		return nil
	}

	return unsetCloudCredentials
}

type UnsetCloudCredentials func(ICloudGetter[iota.Cloud]) error
