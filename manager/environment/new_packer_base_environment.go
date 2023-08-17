package environment

import "github.com/spf13/viper"

type PackerBaseRequired struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type PackerBaseOptional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerBaseEnvironment struct {
	Required *PackerBaseRequired
	Optional *PackerBaseOptional
}

func NewPackerBaseEnvironment() *PackerBaseEnvironment {
	return &PackerBaseEnvironment{
		Required: &PackerBaseRequired{
			GIT_USERNAME: viper.GetString("Git.Username"),
			GIT_EMAIL:    viper.GetString("Git.Email"),
			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
		},
		Optional: &PackerBaseOptional{
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("Git.Hub.PersonalAccessTokenClassic"),
		},
	}
}
