package packer

import (
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
)

func NewPackerTemplateCombo(
	options ...Option,
) (
	packerTemplateCombo *PackerTemplateCombo,
) {
	var (
		option Option
	)

	packerTemplateCombo = &PackerTemplateCombo{}
	for _, option = range options {
		option(packerTemplateCombo)
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

	option = func(packerTemplate *PackerTemplateCombo) (err error) {
		packerTemplate.GeneralTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.PACKER,
			constants.GENERAL_DIR_NAME,
			constants.PACKER_GENERAL_VARS_TEMPLATE,
		)

		switch cloud.Kind {
		case constants.Aws:
			packerTemplate.CloudTemplateAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.PACKER,
				constants.AWS,
				constants.PACKER_AWS_VARS_TEMPLATE,
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

type PackerTemplateCombo struct {
	GeneralTemplateAbsPath string
	CloudTemplateAbsPath   string
}

type Option func(*PackerTemplateCombo) error
