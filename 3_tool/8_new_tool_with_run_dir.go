package tool

import (
	"path/filepath"

	constants "github.com/ed3899/kumo/0_constants"
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/samber/oops"
)

type ToolWithRunDir struct {
	ToolWithExecAbsPath
	RunDir string
}

func NewToolWithRunDir(
	toolKind constants.ToolKind,
	toolWithExecAbsPath ToolWithExecAbsPath,
	cloud cloud.Cloud,
	kumoExecutableAbsPath string,
) (toolWithRunDir ToolWithRunDir, err error) {

	var (
		oopsBuilder = oops.Code("new_tool_setup_failed").
			With("tool", toolKind).
			With("cloud", cloud).
			With("kumoExecutableAbsPath", kumoExecutableAbsPath).
			With("toolWithExecAbsPath", toolWithExecAbsPath)
	)

	switch toolKind {
	case constants.Packer:
		toolWithRunDir = ToolWithRunDir{
			ToolWithExecAbsPath: toolWithExecAbsPath,
			RunDir: filepath.Join(
				kumoExecutableAbsPath,
				constants.PACKER,
				cloud.Name,
			),
		}

	case constants.Terraform:
		toolWithRunDir = ToolWithRunDir{
			ToolWithExecAbsPath: toolWithExecAbsPath,
			RunDir: filepath.Join(
				kumoExecutableAbsPath,
				constants.TERRAFORM,
				cloud.Name,
			),
		}

	default:
		err = oopsBuilder.
			Errorf("Unknown tool kind: %d", toolKind)
	}

	return
}
