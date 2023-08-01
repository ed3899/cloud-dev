package general

type Environment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	switch {
	case e.GIT_USERNAME == "":
		return true
	case e.GIT_EMAIL == "":
		return true
	case len(e.ANSIBLE_TAGS) == 0:
		return true
	default:
		return false
	}
}
