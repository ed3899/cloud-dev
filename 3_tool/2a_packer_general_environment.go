package tool

type PackerGeneralRequired struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type PackerGeneralOptional struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerGeneralEnvironment struct {
	Required  PackerGeneralRequired
	Optional PackerGeneralOptional
}

func (pge PackerGeneralEnvironment) IsEnvironment() bool {
	return true
}

func (pge PackerGeneralEnvironment) IsGeneralEnvironment() bool {
	return true
}