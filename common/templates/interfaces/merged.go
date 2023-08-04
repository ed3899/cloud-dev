package interfaces

import "text/template"

type Merged interface {
	AbsPath() (path string)
	Name() (name string)
	Instance() (instance *template.Template)
}