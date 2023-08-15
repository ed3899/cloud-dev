package manager

import (
	"github.com/ed3899/kumo/utils/file"
)

func CheckIfManagerPathExecutableExists(
	manager *Manager,
) bool {
	return file.IsFilePresent(manager.Path.Executable)
}
