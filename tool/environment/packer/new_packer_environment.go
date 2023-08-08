package packer

import (
	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool/environment/packer/general"
	"github.com/ed3899/kumo/utils/environment"
	"github.com/samber/oops"
)

func NewPackerEnvironment(
	options ...Option,
) (
	packerEnvironment *PackerEnvironment,
	err error,
) {

	var (
		oopsBuilder = oops.
				Code("NewPackerEnvironment").
				With("options", options)

		opt Option
	)

	packerEnvironment = &PackerEnvironment{}
	for _, opt = range options {
		if err = opt(packerEnvironment); err != nil {
			err = oopsBuilder.
				With("packerEnvironment", packerEnvironment).
				Wrapf(err, "Error while applying option %v", opt)
			return
		}
	}

	return

}

func WithPackerGeneralEnvironment(
	packerGeneralEnvironment *general.PackerGeneralEnvironment,
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithGeneralEnvironment").
				With("packerGeneralEnvironment", packerGeneralEnvironment).
				With("areRequiredFieldsNotFilled", areRequiredFieldsNotFilled)

		requiredFieldsNotFilled bool
		missingField            string
	)

	option = func(packerEnvironment *PackerEnvironment) (err error) {
		requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(packerGeneralEnvironment.Required)
		if requiredFieldsNotFilled {
			err = oopsBuilder.
				With("packerGeneralEnvironment.Required", packerGeneralEnvironment.Required).
				Errorf("Required field '%s' is not filled", missingField)
			return
		}

		packerEnvironment.General = packerGeneralEnvironment

		return
	}
	return

}

func WithCloudChoice(
	cloud cloud.Cloud,
	packerCloudEnvironmentOptions *PackerCloudEnvironmentOptions,
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
				Code("WithCloudEnvironment").
				With("packerCloudEnvironment", packerCloudEnvironmentOptions).
				With("areRequiredFieldsNotFilled", areRequiredFieldsNotFilled)

		requiredFieldsNotFilled bool
		missingField            string
	)

	option = func(packerEnvironment *PackerEnvironment) (err error) {
		switch cloud.Kind {
		case constants.Aws:
			requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(packerCloudEnvironmentOptions.Aws.Required)
			if requiredFieldsNotFilled {
				err = oopsBuilder.
					With("packerCloudEnvironment.Aws.Required", packerCloudEnvironmentOptions.Aws.Required).
					Errorf("Required field '%s' is not filled", missingField)
				return
			}

			packerEnvironment.Cloud = packerCloudEnvironmentOptions.Aws

		default:
			err = oopsBuilder.
				With("cloud.Kind", cloud.Kind).
				Errorf("Cloud '%s' is not supported", cloud.Name)
			return
		}

		return
	}

	return
}

type PackerEnvironment struct {
	General *general.PackerGeneralEnvironment
	Cloud   interfaces.PackerCloudEnvironmentI
}

type Option func(*PackerEnvironment) (err error)
