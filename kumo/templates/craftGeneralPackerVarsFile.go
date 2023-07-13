package templates

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type GeneralPackerEnvironment struct {
	*AWS_PackerEnvironment
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

func CraftGeneralPackerVarsFile(cloud string) (generalPackerVarsPath string, err error) {
	genEnv := &GeneralPackerEnvironment{
		GIT_USERNAME:                          viper.GetString("Git.Username"),
		GIT_EMAIL:                             viper.GetString("Git.Email"),
		ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
		GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
	}

	generalPackerVarsPath, err = CraftGenericPackerVarsFile[*GeneralPackerEnvironment](cloud,"GeneralPackerVarsTemplate.tmpl", "general_ami.auto.pkrvars.hcl", genEnv)

	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting generic Packer Vars file")
		return "", err
	}

	return generalPackerVarsPath, nil

}
