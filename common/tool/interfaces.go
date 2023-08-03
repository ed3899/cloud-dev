package tool

type ConfigI interface {
	GetDependenciesDirName() (dependenciesDirName string)
	GetType() (toolType Type)
	GetName() (toolName string)
	GetZipName() (toolZipName string)
	GetZipAbsPath() (toolZipAbsPath string)
	GetZipContentLength() (toolZipContentLength int64, err error)
	GetExecutableName() (toolExecutableName string)
	GetVersion() (toolVersion string)
	GetInitialDir() (initialDir string)
	GetDir() (toolDir string)
	GetUrl() (toolUrl string)
	GoInitialDir() (err error)
	GoDir() (err error)
}
