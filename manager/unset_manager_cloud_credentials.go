package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func (m *Manager) UnsetManagerCloudCredentials() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetCloudCredentials")

	switch m.Cloud {
	case iota.Aws:
		for key := range awsCredentials {
			if err := os.Unsetenv(key); err != nil {
				return oopsBuilder.
					With("cloudName", m.Cloud.Name()).
					Wrapf(err, "failed to unset environment variable %s", key)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", m.Cloud.Name()).
			Errorf("unknown cloud: %#v", m.Cloud)
	}

	return nil

}
