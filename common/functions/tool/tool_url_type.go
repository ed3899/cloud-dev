package tool

import (
	"github.com/ed3899/kumo/common/functions/host"
)

type ToolUrlWith func(*ToolUrlWithArgs[
	ToolName,
	ToolVersion,
	host.CurrentOs,
	host.CurrentArch,
]) ToolUrl

type ToolUrl func() string

type ToolUrlWithArgs[
	TN ToolName,
	TV ToolVersion,
	CO host.CurrentOs,
	CA host.CurrentArch,
] struct {
	ToolName    TN
	ToolVersion TV
	CurrentOs   CO
	CurrentArch CA
}
