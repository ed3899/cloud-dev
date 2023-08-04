package structs

import "github.com/ed3899/kumo/common/templates/interfaces"

type MergedEnvironment[E interfaces.Environment] struct {
	General E
	Cloud   E
}
