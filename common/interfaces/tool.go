package interfaces

import (
	"github.com/ed3899/kumo/common/alias"
	"github.com/ed3899/kumo/common/constants"
)

type ToolI interface {
	KindGetter[constants.Tool]
	NameGetter[alias.ToolName]
	VersionGetter[alias.ToolVersion]
	UrlGetter[alias.ToolUrl]
}
