package templates

import (
	"os"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

// Crafts a generic cloud Terraform Vars file.
//
// It looks for templates under "./templates/terraform/<cloud>/<templateName>.tmpl"
//
// It creates a file under "./terraform/<cloud>/<terraformVarsFileName>.auto.tfvars"
//
// The env parameter is a pointer to a struct that contains the data
// to be used in the template. Example:
//
//	awsEnv := &AWS_TerraformEnvironment{
//		AWS_REGION:                   viper.GetString("AWS.Region"),
//		AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
//		AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
//		AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
//	}
//
// Usage:
//
//	("aws", "AWS_TerraformTfVars.tmpl", "aws.auto.tfvars", awsEnv) -> ("terraform/aws/aws.auto.tfvars", nil)
func CraftGenericCloudTerraformTfVarsFile(cloud, templateName, terraformVarsFileName string, env any) (resultingTerraformTfVarsPath string, err error) {
	// Get template
	templatePath, err := utils.CraftAbsolutePath("templates", "terraform", cloud, templateName)
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
