package templates

import (
	"os"
	"text/template"

	templates_packer_aws "github.com/ed3899/kumo/templates/packer/aws"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type GeneralPackerEnvironment struct {
	*templates_packer_aws.AWS_PackerEnvironment
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

	// Get template
	generalTmplName := "GeneralPackerVars.tmpl"
	templatePath, err := utils.CraftAbsolutePath("templates", "packer", "general", generalTmplName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", generalTmplName)
		return "", err
	}

	// Create vars file
	generalPackerVarsFileName := "general.auto.tfvars"
	generalPackerVarsPath, err = utils.CraftAbsolutePath("terraform", cloud, generalPackerVarsFileName)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s file", generalPackerVarsFileName)
		return "", err
	}

	file, err := os.Create(generalPackerVarsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", generalPackerVarsFileName)
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, genEnv)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing %s template file", generalTmplName)
		return "", err
	}

	// Return path to vars file
	return generalPackerVarsPath, nil
}
