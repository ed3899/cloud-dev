package general

import (
	"github.com/ed3899/kumo/utils/environment"
	"github.com/samber/oops"
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

func NewPackerGeneralEnvironment(
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
	options ...Option,
) (packerGeneralEnvironment *PackerGeneralEnvironment, err error) {
	var (
		oopsBuilder = oops.
				Code("NewPackerGeneralEnvironment").
				With("opts", options).
				With("requiredFieldsAreNotFilled", areRequiredFieldsNotFilled)

		opt                     Option
		requiredFieldsNotFilled bool
		missingField            string
	)

	packerGeneralEnvironment = &PackerGeneralEnvironment{}
	for _, opt = range options {
		opt(packerGeneralEnvironment)
	}

	requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(packerGeneralEnvironment.Required)
	if requiredFieldsNotFilled {
		err = oopsBuilder.
			With("packerGeneralEnvironment.Required", packerGeneralEnvironment.Required).
			Errorf("Required field '%s' is not filled", missingField)
		return
	}

	return
}

func WithGitUsername(gitUsername string) (option Option) {
	option = func(environment *PackerGeneralEnvironment) {
		environment.Required.GIT_USERNAME = gitUsername
	}

	return
}

func WithGitEmail(gitEmail string) (option Option) {
	option = func(environment *PackerGeneralEnvironment) {
		environment.Required.GIT_EMAIL = gitEmail
	}

	return
}

func WithAnsibleTags(ansibleTags []string) (option Option) {
	option = func(environment *PackerGeneralEnvironment) {
		environment.Required.ANSIBLE_TAGS = ansibleTags
	}

	return
}

func WithGitHubPersonalAccessTokenClassic(gitHubPersonalAccessTokenClassic string) (option Option) {
	option = func(environment *PackerGeneralEnvironment) {
		environment.Optional.GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC = gitHubPersonalAccessTokenClassic
	}

	return
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
	Required *PackerGeneralRequired
	Optional *PackerGeneralOptional
}

type Option func(environment *PackerGeneralEnvironment)

type NewPackerGeneralEnvironmentF func(
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
	options ...Option,
) (packerGeneralEnvironment *PackerGeneralEnvironment, err error)
