package templates

type EnvironmentI interface {
	IsNotValidEnvironment() bool
}

type TemplateSingle interface {
	GetParentDirName() string
	GetName() string
	GetEnvironment() EnvironmentI
}
