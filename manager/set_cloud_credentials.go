package manager

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func SetCloudCredentialsWith(
	osSetenv func(string, string) error,
) SetCloudCredentials {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCloudCredentials")

	setCloudCredentials := func(manager ICloudGetter[iota.Cloud]) error {
		switch manager.Cloud() {
		case iota.Aws:
			for key, value := range awsCredentials {
				if err := osSetenv(key, value); err != nil {
					return oopsBuilder.
						With("cloudName", manager.Cloud().Name()).
						Wrapf(err, "failed to set environment variable %s to %s", key, value)
				}
			}

		default:
			return oopsBuilder.
				With("cloudName", manager.Cloud().Name()).
				Errorf("unknown cloud: %#v", manager.Cloud())
		}

		return nil
	}

	return setCloudCredentials
}

type SetCloudCredentials func(ICloudGetter[iota.Cloud]) error
