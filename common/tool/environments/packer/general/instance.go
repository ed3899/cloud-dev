package general

import "github.com/ed3899/kumo/common/utils"

type Required struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type Optional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type Environment struct {
	Required  *Required
	Optional *Optional
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructCompletellyFilled(e.Required)
}
