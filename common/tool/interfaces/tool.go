package interfaces

import (
	"github.com/ed3899/kumo/common/tool/constants"
)

type Tool interface {
	Kind() (toolKind constants.Kind)
	Name() (toolName string)
	SetPluginPath() (err error)
	UnsetPluginPath() (err error)
	ExecutableName() (toolExecutableName string)
	Version() (toolVersion string)
	RunDir() (toolDir string)
	Url() (toolUrl string)
}
