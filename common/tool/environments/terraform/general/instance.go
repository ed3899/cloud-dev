package general

import "github.com/ed3899/kumo/common/utils"

type Required struct {
	ALLOWED_IP string
}

type Environment struct {
	Required *Required
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructCompletellyFilled(e.Required)
}
