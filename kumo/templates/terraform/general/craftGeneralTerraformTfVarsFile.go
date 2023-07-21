package templates

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	templates_terraform_utils "github.com/ed3899/kumo/templates/terraform/utils"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type GeneralTerraformEnvironment struct {
	*templates_terraform_aws.AWS_TerraformEnvironment
	AMI_ID     string
	ALLOWED_IP string
}

func CraftGeneralTerraformTfVarsFile(gte *GeneralTerraformEnvironment) (generalTerraformVarsPath string, err error) {
	// Get template
	generalTmplName := "GeneralTerraformTfVars.tmpl"
	templatePath, err := filepath.Abs(filepath.Join("templates", "terraform", "general", generalTmplName))
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", generalTmplName)
		return "", err
	}

	// Create vars file
	generalTerraformVarsFileName := "general.auto.tfvars"
	generalTerraformVarsPath, err = filepath.Abs(filepath.Join("terraform", cloud, generalTerraformVarsFileName))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s file", generalTerraformVarsFileName)
		return "", err
	}

	file, err := os.Create(generalTerraformVarsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", generalTerraformVarsFileName)
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, *gte)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing %s template file", generalTmplName)
		return "", err
	}

	// Return path to vars file
	return generalTerraformVarsPath, nil

}
