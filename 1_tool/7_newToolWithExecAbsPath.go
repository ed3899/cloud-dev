package tool

import (
	"fmt"
	"os"
	"path/filepath"

	constants "github.com/ed3899/kumo/0_constants"
	"github.com/samber/oops"
)

type ToolWithExecAbsPath struct {
	Tool
	ExecutableAbsPath string
}

func NewToolWithExecAbsPath(tool Tool) (toolWithExecAbsPath ToolWithExecAbsPath, err error) {
	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed")

		kumoExecAbsPath string
	)

	if kumoExecAbsPath, err = os.Executable(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while retrieving absolute path to kumo executable")
		return
	}

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
