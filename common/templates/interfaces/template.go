package interfaces

import "text/template"

type MergedTemplate interface {
	AbsPath() (path string)
	Name() (name string)
	Instance() (instance *template.Template)
}

type Template interface {
	ParentDirName() string
	MergedTemplate
}
