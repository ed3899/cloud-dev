package tool

import (
	"path/filepath"

	constants "github.com/ed3899/kumo/0_constants"
	cloud "github.com/ed3899/kumo/2_cloud"
	"github.com/samber/oops"
)

type ToolWithPluginDir struct {
	ToolWithRunDir
	PluginsDir string
}

func NewToolWithPluginDir(
	toolKind constants.ToolKind,
	toolWithRunDir ToolWithRunDir,
	cloud cloud.Cloud,
	kumoExecutableAbsPath string,
) (toolWithPluginDir ToolWithPluginDir, err error) {
	var (
		oopsBuilder = oops.Code(
			"new_tool_with_plugin_dir_failed",
		)
	)

	switch toolKind {
	case constants.Packer:
		toolWithPluginDir = ToolWithPluginDir{
			ToolWithRunDir: toolWithRunDir,
			PluginsDir: filepath.Join(
				kumoExecutableAbsPath,
				constants.PACKER,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			),
		}

	case constants.Terraform:
		toolWithPluginDir = ToolWithPluginDir{
			ToolWithRunDir: toolWithRunDir,
			PluginsDir: filepath.Join(
				kumoExecutableAbsPath,
				constants.TERRAFORM,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			),
		}

	default:
		err = oopsBuilder.
			Errorf("Unknown tool kind: %d", toolKind)
		return
	}

	return
}
