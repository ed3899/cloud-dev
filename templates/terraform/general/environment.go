package general

type Environment struct {
	ALLOWED_IP string
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return e.ALLOWED_IP == ""
}
