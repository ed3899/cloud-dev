package interfaces

import (
	"text/template"

	"github.com/ed3899/kumo/common/templates"
)

type Aws interface {
	AbsPath() (absPath string)
	ParentDirName() (dir string)
	Instance() (instance *template.Template)
	Environment() (environment templates.EnvironmentI)
}
