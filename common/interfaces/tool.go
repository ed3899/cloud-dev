package interfaces

import "github.com/ed3899/kumo/constants"

type KindGetter[K ~int] interface {
	Kind() K
}

type NameGetter interface {
	Name() string
}

type VersionGetter interface {
	Version() string
}

type UrlGetter interface {
	Url() string
}

type AbsPathGetter interface {
	AbsPath() *AbsPath
}

type ToolConfig interface {
	KindGetter[constants.ToolKind]
	NameGetter
	VersionGetter
	UrlGetter
	AbsPathGetter
}

type ExecutableGetter interface {
	Executable() string
}

type DirGetter interface {
	Dir() string
}

type TemplateFileGetter interface {
	TemplateFile() *TemplateFileCombo
}

type AbsPath interface {
	ExecutableGetter
	DirGetter
	TemplateFileGetter
}

type CloudGetter interface {
	Cloud() string
}

type GeneralGetter interface {
	General() string
}

type TemplateFileCombo interface {
	CloudGetter
	GeneralGetter
}
