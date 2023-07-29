package general

import (
	"github.com/ed3899/kumo/utils"
)

type Environment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructFilledCompletely(e)
}
