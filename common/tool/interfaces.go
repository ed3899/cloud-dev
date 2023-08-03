package tool

type ConfigI interface {
	GetDependenciesDirName() (dependenciesDirName string)
	GetToolType() (toolType ToolType)
	GetToolName() (toolName string)
	GetToolVersion() (toolVersion string)
	GetInitialDir() (initialDir string)
	GetToolDir() (toolDir string)
	GoInitialDir() (err error)
	GoTargetDir() (err error)
}
