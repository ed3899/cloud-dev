package manager

import (
	"github.com/ed3899/kumo/utils/file"
)

// Return true if the tool executable exists
func (m *Manager) ToolExecutableExists() bool {
	return file.IsFilePresent(m.Path.Executable)
}
