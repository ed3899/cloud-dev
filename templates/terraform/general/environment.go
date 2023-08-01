package general

import (
	"github.com/ed3899/kumo/common/utils"
)

type Environment struct {
	ALLOWED_IP string
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructCompletellyFilled(e)
}
