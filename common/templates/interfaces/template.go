package interfaces

import "text/template"

type Template interface {
	ParentDirName() string
	AbsPath() string
	Environment() Environment
	Instance() (instance *template.Template)
}
