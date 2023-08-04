package interfaces

import (
	"github.com/ed3899/kumo/common/cloud/constants"
)

type Cloud interface {
	Name() (cloudName string)
	Type() (cloudType constants.Kind)
}
