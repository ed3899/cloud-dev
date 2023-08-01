package templates

type EnvironmentI interface {
	IsNotValidEnvironment() bool
}

type TemplateSingle interface {
	GetParentDirName() string
	GetAbsPath() string
	GetEnvironment() EnvironmentI
}

type PackerManifestI interface {
	GetLastBuiltAmiId() (lastBuiltAmiId string)
}
