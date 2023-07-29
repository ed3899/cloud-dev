package templates

type TemplateSingle interface {
	GetParentDirName() string
	GetName() string
	GetEnvironment() Environment
}