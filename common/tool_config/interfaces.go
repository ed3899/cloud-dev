package tool_config

type ToolI interface {
	DependenciesDirName() (dependenciesDirName string)
	Kind() (toolKind Kind)
	Name() (toolName string)
	ZipAbsPath() (toolZipAbsPath string)
	ZipContentLength() (toolZipContentLength int64, err error)
	ExecutableName() (toolExecutableName string)
	Version() (toolVersion string)
	InitialDir() (initialDir string)
	Dir() (toolDir string)
	Url() (toolUrl string)
	GoInitialDir() (err error)
	GoDir() (err error)
}
