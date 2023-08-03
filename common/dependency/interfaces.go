package dependency

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