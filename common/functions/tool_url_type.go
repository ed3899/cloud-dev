package functions

type ToolUrlWith func(*ToolUrlWithArgs[func() string]) ToolUrl

type ToolUrl func() string

type ToolUrlWithArgs[F any] struct {
	ToolName    F
	ToolVersion F
	CurrentOs   F
	CurrentArch F
}
