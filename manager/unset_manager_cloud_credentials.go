package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func UnsetManagerCloudCredentials(
	manager *Manager,
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCloudCredentials")

	switch manager.Cloud {
	case iota.Aws:
		for key := range awsCredentials {
			if err := os.Unsetenv(key); err != nil {
				return oopsBuilder.
					With("cloudName", manager.Cloud.Name).
					Wrapf(err, "failed to unset environment variable %s", key)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", manager.Cloud.Name).
			Errorf("unknown cloud: %#v", manager.Cloud)
	}

	return nil

}
