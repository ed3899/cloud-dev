package general

import (
	"github.com/ed3899/kumo/utils/environment"
)

// func NewEnvironment() (generalEnvironment PackerGeneralEnvironment) {
// 	generalEnvironment = PackerGeneralEnvironment{
// 		Required: PackerGeneralRequired{
// 			GIT_USERNAME: viper.GetString("Git.Username"),
// 			GIT_EMAIL:    viper.GetString("Git.Email"),
// 			ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
// 		},
// 		Optional: PackerGeneralOptional{
// 			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
// 		},
// 	}

// 	return
// }

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

type Option func(environment *PackerGeneralEnvironment)

type NewPackerGeneralEnvironmentF func(
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
	options ...Option,
) (packerGeneralEnvironment *PackerGeneralEnvironment, err error)
