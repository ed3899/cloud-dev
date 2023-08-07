package general

import (
	"github.com/spf13/viper"
)

func NewEnvironment() (generalEnvironment Environment) {
	generalEnvironment = Environment{
		Required: Required{
			GIT_USERNAME: viper.GetString("Git.Username"),
			GIT_EMAIL:    viper.GetString("Git.Email"),
			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
		},
		Optional: Optional{
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		},
	}

	return
}

type Required struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type Optional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type Environment struct {
	Required Required
	Optional Optional
}


type NewEnvironmentF func() Environment
