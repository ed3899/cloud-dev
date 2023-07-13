package templates

import (
	"os"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type PackerEnvironmentI interface {
	*AWS_PackerEnvironment | *GeneralPackerEnvironment
}

func CraftGenericPackerVarsFile[E PackerEnvironmentI](cloud, templateName, packerVarsFileName string, env E) (resultingPackerVarsPath string, err error) {
	// Get template
	templatePath, err := utils.CraftAbsolutePath("templates", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to Packer AWS Vars template file")
		return "", err
	}

	// Create vars file
	resultingPackerVarsPath, err = utils.CraftAbsolutePath("packer", cloud, packerVarsFileName)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to Packer AWS Vars file")
		return "", err
	}

	file, err := os.Create(resultingPackerVarsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, env)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}

	// Return path to vars file
	return resultingPackerVarsPath, nil
}
