package aws

import (
	"html/template"
	"path/filepath"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func LoadPackerAwsVarsTemplate(
	kumoExecAbsPath string,
) (
	packerAwsTemplate *template.Template,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("LoadPackerAwsVarsTemplate").
				With("kumoExecAbsPath", kumoExecAbsPath)

		packerAwsTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.PACKER,
			constants.AWS,
			constants.PACKER_AWS_VARS_TEMPLATE,
		)
	)

	if packerAwsTemplate, err = template.ParseFiles(packerAwsTemplateAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse packer AWS template '%s'", packerAwsTemplateAbsPath)
		return
	}

	return

}
