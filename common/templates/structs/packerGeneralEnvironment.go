package structs

import "github.com/ed3899/kumo/common/utils"

type PackerGeneralRequired struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type PackerGeneralOptional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerGeneralEnvironment struct {
	Required  *PackerGeneralRequired
	Optional *PackerGeneralOptional
}

func (pge *PackerGeneralEnvironment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructCompletellyFilled(pge.Required)
}
