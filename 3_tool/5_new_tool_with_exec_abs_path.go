package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/constants"
)

type ToolWithExecAbsPath struct {
	Tool
	ExecutableAbsPath string
}

func NewToolWithExecAbsPath(tool Tool, kumoExecAbsPath string) (toolWithExecAbsPath ToolWithExecAbsPath) {
	toolWithExecAbsPath = ToolWithExecAbsPath{
		Tool: tool,
		ExecutableAbsPath: filepath.Join(
			kumoExecAbsPath,
			constants.DEPENDENCIES_DIR_NAME,
			tool.Name,
			fmt.Sprintf("%s.exe", tool.Name),
		),
	}

	return
}
