package structs

import "github.com/ed3899/kumo/common/templates"

type Environment[E templates.EnvironmentI] struct {
	General E
	Cloud   E
}
