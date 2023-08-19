package manager

import (
	"os"

	"github.com/samber/oops"
)

func (m *Manager) SetPluginsEnvironmentVars() error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetPluginsEnvironmentVars")

	err := os.Setenv(m.Tool.PluginPathEnvironmentVariable(), m.Path.Dir.Plugins)
	if err != nil {
		return oopsBuilder.
			Errorf(
				"failed to set %s environment variable to %s",
				m.Tool.PluginPathEnvironmentVariable(),
				m.Path.Dir.Plugins,
			)
	}

	return nil
}
