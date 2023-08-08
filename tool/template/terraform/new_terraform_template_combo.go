package terraform

import (
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func NewTerraformTemplateCombo(
	options ...Option,
) (
	terraformTemplateCombo *TerraformTemplateCombo,
) {
	var (
		option Option
	)

	terraformTemplateCombo = &TerraformTemplateCombo{}
	for _, option = range options {
		option(terraformTemplateCombo)
	}

	return
}

func WithCloudChoice(
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithCloudChoice").
			With("cloud", cloud).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(packerTemplate *TerraformTemplateCombo) (err error) {
		packerTemplate.GeneralTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.TERRAFORM,
			constants.GENERAL_DIR_NAME,
			constants.TERRAFORM_GENERAL_VARS_TEMPLATE,
		)

		switch cloud.Kind {
		case constants.Aws:
			packerTemplate.CloudTemplateAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.TERRAFORM,
				constants.AWS,
				constants.TERRAFORM_AWS_VARS_TEMPLATE,
			)
		default:
			err = oopsBuilder.
				Wrapf(err, "Unsupported cloud '%s'", cloud.Name)
			return
		}

		return
	}

	return
}

type TerraformTemplateCombo struct {
	GeneralTemplateAbsPath string
	CloudTemplateAbsPath   string
}

type Option func(*TerraformTemplateCombo) error
