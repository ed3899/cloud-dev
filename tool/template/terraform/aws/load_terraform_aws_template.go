package terraform

import (
	"html/template"
	"path/filepath"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func LoadTerraformAwsVarsTemplate(
	kumoExecAbsPath string,
) (
	terraformAwsTemplate *template.Template,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("LoadTerraformAwsVarsTemplate").
				With("kumoExecAbsPath", kumoExecAbsPath)

		terraformAwsTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.TERRAFORM,
			constants.AWS,
			constants.TERRAFORM_AWS_VARS_TEMPLATE,
		)
	)

	if terraformAwsTemplate, err = template.ParseFiles(terraformAwsTemplateAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse terraform AWS template '%s'", terraformAwsTemplateAbsPath)
		return
	}

	return
}
