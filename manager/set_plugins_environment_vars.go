package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) SetPluginsEnvironmentVars() error {
	oopsBuilder := oops.
		Code("SetPluginsEnvironmentVars")

	err := os.Setenv(m.Tool.PluginPathEnvironmentVariable(), m.Path.Plugins)
	if err != nil {
		return oopsBuilder.
			Errorf("failed to set %s environment variable to %s", m.Tool.PluginPathEnvironmentVariable(), m.Path.Plugins)
	}

	return nil
}
