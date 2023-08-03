package packer

import (
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
)

type Main interface {
	GetDependenciesDirName() string
	GetToolName() string
	GetToolVersion() string
	GetToolExecutableName() string
	GetToolZipName() string
}

type Utils interface {
	GetKumoExecutableAbsPath() (absPathToKumoExecutable string, err error)
	GetCurrentHostSpecs() (currentOs string, currentArch string)
	CreateHashicorpURL(toolName string, toolVersion string, currentOs string, currentArch string) (url string)
	GetContentLength(url string) (contentLength int64, err error)
}

type ConfigI interface {
	Main
	Utils
}

type Config struct {
	dependenciesDirName string
	toolName            string
	toolVersion         string
}

func NewConfig(toolType tool.ToolType) (config *Config, err error) {
	var (
		oopsBuilder = oops.
				Code("new_config_failed").
				With("toolType", toolType)

		dependenciesDirName = dirs.DEPENDENCIES_DIR_NAME
	)

	switch toolType {
	case tool.Packer:
		config = &Config{
			dependenciesDirName: dependenciesDirName,
			toolName:            tool.PACKER_NAME,
			toolVersion:         tool.PACKER_VERSION,
		}

	case tool.Terraform:
		config = &Config{
			dependenciesDirName: dependenciesDirName,
			toolName:            tool.TERRAFORM_NAME,
			toolVersion:         tool.TERRAFORM_VERSION,
		}

	default:
		err = oopsBuilder.
			With("toolType", toolType).
			Errorf("invalid tool type")
		return
	}

	return
}

func (c *Config) GetKumoExecutableAbsPath() (absPathToKumoExecutable string, err error) {
	return
}

func (c *Config) GetCurrentHostSpecs() (currentOs string, currentArch string) {
	return
}

func (c *Config) CreateHashicorpURL(toolName string, toolVersion string, currentOs string, currentArch string) (url string) {
	return
}

func (c *Config) GetContentLength(url string) (contentLength int64, err error) {
	return
}
