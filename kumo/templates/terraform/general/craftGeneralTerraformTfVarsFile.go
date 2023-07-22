package templates

import (
	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type GeneralTerraformEnvironment struct {
	*templates_terraform_aws.AWS_TerraformEnvironment
	AMI_ID     string
	ALLOWED_IP string
}

func CraftGeneralTerraformTfVarsFile(gte *GeneralTerraformEnvironment) (generalTerraformVarsPath string, err error) {
	// Get template
	tmpl, err := utils.CreateTemplate(&utils.TemplateProps{
		Tool: utils.Terraform,
		Kind: utils.General,
	})
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while getting template for terraform")
		return "", err
	}

	// Create vars file
	varsFile, err := utils.CreateVarsFile(&utils.TemplateProps{
		Tool: utils.Terraform,
		Kind: utils.General,
	})
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating vars file")
		return "", err
	}

	err = utils.ExecuteTemplate(&utils.ExecuteTemplateProps{
		Template:        tmpl,
		VarsFileAbsPath: varsFile.Name(),
		Environment:     gte,
	})
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing template")
		return "", err
	}

	generalTerraformVarsPath = varsFile.Name()

	return
}
