package manager

import (
	"github.com/ed3899/kumo/utils/file"
)

func (m *Manager) ToolExecutableExists() bool {
	return file.IsFilePresent(m.Path.Executable)
}
