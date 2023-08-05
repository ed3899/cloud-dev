package environment

import (
	"github.com/ed3899/kumo/constants"
	"github.com/spf13/viper"
)

func NewPackerGeneralEnvironment(toolKind constants.ToolKind) (generalEnvironment GeneralEnvironmentI) {
	return PackerGeneralEnvironment{
		Required: PackerGeneralRequired{
			GIT_USERNAME: viper.GetString("Git.Username"),
			GIT_EMAIL:    viper.GetString("Git.Email"),
			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
		},
		Optional: PackerGeneralOptional{
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		},
	}

}

func (pge PackerGeneralEnvironment) IsEnvironment() bool {
	return true
}

func (pge PackerGeneralEnvironment) IsGeneralEnvironment() bool {
	return true
}

type PackerGeneralRequired struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type PackerGeneralOptional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerGeneralEnvironment struct {
	Required PackerGeneralRequired
	Optional PackerGeneralOptional
}
