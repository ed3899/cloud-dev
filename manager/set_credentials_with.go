package manager

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func SetCredentialsWith(
	osSetenv func(string, string) error,
) ForCloudGetter {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCredentialsWith")

	forManager := func(manager ICloudGetter) error {
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

	return forManager
}
