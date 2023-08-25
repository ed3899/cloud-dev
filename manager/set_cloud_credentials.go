package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

// Sets the cloud credentials environment variables.
func (m *Manager) SetCloudCredentials() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCloudCredentials")

	switch m.Cloud {
	case iota.Aws:
		for key, value := range awsCredentials {
			if err := os.Setenv(key, viper.GetString(value)); err != nil {
				return oopsBuilder.
					With("cloudName", m.Cloud.Name).
					Wrapf(err, "failed to set environment variable %s to %s", key, value)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", m.Cloud.Name).
			Errorf("unknown cloud: %#v", m.Cloud)
	}

	return nil
}
