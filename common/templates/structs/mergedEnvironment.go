package structs

import "github.com/ed3899/kumo/common/templates/interfaces"

type Environment[E interfaces.Environment] struct {
	General E
	Cloud   E
}