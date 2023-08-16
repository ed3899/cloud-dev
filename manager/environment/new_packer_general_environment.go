package environment

import "github.com/spf13/viper"

type PackerGeneralRequired struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type PackerGeneralOptional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerGeneralEnvironment struct {
	Required *PackerGeneralRequired
	Optional *PackerGeneralOptional
}

func NewPackerGeneralEnvironment() *PackerGeneralEnvironment {
	return &PackerGeneralEnvironment{
		Required: &PackerGeneralRequired{
			GIT_USERNAME: viper.GetString("Git.Username"),
			GIT_EMAIL:    viper.GetString("Git.Email"),
			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
		},
		Optional: &PackerGeneralOptional{
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("Git.Hub.PersonalAccessTokenClassic"),
		},
	}
}
