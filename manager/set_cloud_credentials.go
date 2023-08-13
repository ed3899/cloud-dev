package manager

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func SetCloudCredentials(
	osSetenv func(string, string) error,
	manager ICloudGetter[iota.Cloud],
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCloudCredentials")

	managerCloudName := manager.Cloud().Name()

	switch manager.Cloud() {
	case iota.Aws:
		for key, value := range awsCredentials {
			if err := osSetenv(key, value); err != nil {
				return oopsBuilder.
					With("cloudName", managerCloudName).
					Wrapf(err, "failed to set environment variable %s to %s", key, value)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", managerCloudName).
			Errorf("unknown cloud: %#v", manager.Cloud())
	}

	return nil

}
