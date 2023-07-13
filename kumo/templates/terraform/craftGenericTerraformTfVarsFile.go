package templates

import (
	"os"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type TerraformEnvironmentI interface {
	*AWS_TerraformEnvironment | *GeneralTerraformEnvironment
}

func CraftGenericTerraformTfVarsFile[E TerraformEnvironmentI](cloud, templateName, terraformVarsFileName string, env E) (resultingTerraformTfVarsPath string, err error) {
	// Get template
	templatePath, err := utils.CraftAbsolutePath("templates", "terraform", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", templateName)
		return "", err
	}

	// Create vars file
	resultingTerraformTfVarsPath, err = utils.CraftAbsolutePath("terraform", cloud, terraformVarsFileName)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s file", terraformVarsFileName)
		return "", err
	}

	file, err := os.Create(resultingTerraformTfVarsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", terraformVarsFileName)
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, env)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing %s template file", templateName)
		return "", err
	}

	// Return path to vars file
	return resultingTerraformTfVarsPath, nil
}
