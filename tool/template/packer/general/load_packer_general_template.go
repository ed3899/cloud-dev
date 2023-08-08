package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func LoadPackerGeneralVarsTemplate(
	kumoExecAbsPath string,
) (
	packerGeneralTemplate *template.Template,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("LoadPackerGeneralTemplate").
				With("kumoExecAbsPath", kumoExecAbsPath)

		packerGeneralTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.PACKER,
			constants.GENERAL_DIR_NAME,
			constants.PACKER_GENERAL_VARS_TEMPLATE,
		)
	)

	if packerGeneralTemplate, err = template.ParseFiles(packerGeneralTemplateAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse packer general template '%s'", packerGeneralTemplateAbsPath)
		return
	}

	return
}
