package manager

import (
	"os"

	"github.com/samber/oops"
)

// Unsets the plugins path environment variable.
func (m *Manager) UnsetPluginsEnvironmentVars() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("UnsetPluginsEnvironmentVars")

	err := os.Unsetenv(m.Tool.PluginPathEnvironmentVariable())
	if err != nil {
		return oopsBuilder.
			Errorf("failed to unset %s environment variable", m.Tool.PluginPathEnvironmentVariable())
	}

	return nil
}
