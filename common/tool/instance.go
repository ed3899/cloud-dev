package tool

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/samber/oops"
)

type Config struct {
	toolType            ToolType
	toolName            string
	toolVersion         string
	dependenciesDirName string
	initialDir          string
	toolDir             string
}

func NewConfig(toolType ToolType, cloudSetup cloud.CloudSetupI) (toolConfig *Config, err error) {
	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed").
				With("tool", toolType)

		cwd string
	)

	if cwd, err = os.Getwd(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch toolType {
	case Packer:
		toolConfig = &Config{
			dependenciesDirName: dirs.DEPENDENCIES_DIR_NAME,
			toolType:            Packer,
			toolName:            PACKER_NAME,
			toolVersion:         PACKER_VERSION,
			initialDir:          cwd,
			toolDir:             filepath.Join(PACKER_NAME, cloudSetup.GetCloudName()),
		}

	case Terraform:
		toolConfig = &Config{
			dependenciesDirName: dirs.DEPENDENCIES_DIR_NAME,
			toolType:            Terraform,
			toolName:            TERRAFORM_NAME,
			toolVersion:         TERRAFORM_VERSION,
			initialDir:          cwd,
			toolDir:             filepath.Join(TERRAFORM_NAME, cloudSetup.GetCloudName()),
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolType)
		return
	}

	return
}

func (c *Config) GetDependenciesDirName() (dependenciesDirName string) {
	return c.dependenciesDirName
}

func (c *Config) GetToolType() (toolType ToolType) {
	return c.toolType
}

func (c *Config) GetToolName() (toolName string) {
	return c.toolName
}

func (c *Config) GetToolVersion() (toolVersion string) {
	return c.toolVersion
}

func (c *Config) GetInitialDir() (initialDir string) {
	return c.initialDir
}

func (c *Config) GetToolDir() (toolDir string) {
	return c.toolDir
}

func (c *Config) GoInitialDir() (err error) {
	var (
		oopsBuilder = oops.
			Code("go_initial_dir_failed")
	)

	if err = os.Chdir(c.initialDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", c.initialDir)
		return
	}

	return
}

func (c *Config) GoTargetDir() (err error) {
	var (
		oopsBuilder = oops.
			Code("go_target_dir_failed")
	)

	if err = os.Chdir(c.toolDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", c.toolDir)
		return
	}

	return
}
