package tool_config

type ToolI interface {
	Kind() (toolKind Kind)
	Name() (toolName string)
	ExecutableName() (toolExecutableName string)
	Version() (toolVersion string)
	RunDir() (toolDir string)
	Url() (toolUrl string)
}
