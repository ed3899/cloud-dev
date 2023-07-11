package templates

type GeneralPackerEnvironment struct {
	*AWS_PackerEnvironment
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}