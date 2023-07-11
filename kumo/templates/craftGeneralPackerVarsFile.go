package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
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

func CraftGeneralPackerVarsFile() (generalPackerVarsPath string, err error) {
	genEnv := &GeneralPackerEnvironment{
		GIT_USERNAME:                          viper.GetString("Git.Username"),
		GIT_EMAIL:                             viper.GetString("Git.Email"),
		ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
		GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
	}

	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	generalTemplatePath := filepath.Join(cwd, "templates", "GeneralPackerVarsTemplate.tmpl")
	tmpl, err := template.ParseFiles(generalTemplatePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while parsing Packer AWS Vars template file")
		return "", err
	}

	phcldir, err := utils.GetPackerHclDirPath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL directory path")
		return "", err
	}

	generalPackerVarsPath = filepath.Join(phcldir, "general_ami.pkrvars.hcl")

	file, err := os.Create(generalPackerVarsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	err = tmpl.Execute(file, genEnv)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}

	return generalPackerVarsPath, nil

}
