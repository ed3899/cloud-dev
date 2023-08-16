package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) UnsetPluginsEnvironmentVars() error {
	oopsBuilder := oops.
		Code("UnsetPluginsEnvironmentVars")

	err := os.Unsetenv(m.Tool.PluginPathEnvironmentVariable())
	if err != nil {
		return oopsBuilder.
			Errorf("failed to unset %s environment variable", m.Tool.PluginPathEnvironmentVariable())
	}

	return nil
}
