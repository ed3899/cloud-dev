package manager

import (
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/utils/file"
)

func CheckIfManagerPathExecutableExists(
	manager interfaces.IClone[*Manager],
) bool {
	managerClone := manager.Clone()
	return file.IsFilePresent(managerClone.Path.Executable)
}
