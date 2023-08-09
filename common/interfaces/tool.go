package interfaces

import (
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/alias"
)

type ToolI interface {
	KindGetter[constants.ToolKind]
	NameGetter[alias.ToolName]
	VersionGetter[alias.ToolVersion]
	UrlGetter[alias.ToolUrl]
}
