package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func LoadTerraformGeneralVarsTemplate(
	kumoExecAbsPath string,
) (
	terraformGeneralTemplate *template.Template,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("LoadTerraformGeneralVarsTemplate").
				With("kumoExecAbsPath", kumoExecAbsPath)

		terraformGeneralTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.TERRAFORM,
			constants.GENERAL_DIR_NAME,
			constants.TERRAFORM_GENERAL_VARS_TEMPLATE,
		)
	)

	if terraformGeneralTemplate, err = template.ParseFiles(terraformGeneralTemplateAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse terraform general template '%s'", terraformGeneralTemplateAbsPath)
		return
	}

	return
}
