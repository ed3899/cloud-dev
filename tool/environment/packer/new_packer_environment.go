package environment

import (
	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool/environment/packer/aws"
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

func WithGeneralEnvironment(
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

func WithCloudEnvironment(
	packerCloudEnvironment *PackerCloudEnvironmentOptions,
	cloud cloud.Cloud,
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
				Code("WithCloudEnvironment").
				With("packerCloudEnvironment", packerCloudEnvironment).
				With("areRequiredFieldsNotFilled", areRequiredFieldsNotFilled)

		requiredFieldsNotFilled bool
		missingField            string
	)

	option = func(packerEnvironment *PackerEnvironment) (err error) {
		switch cloud.Kind {
		case constants.Aws:
			requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(packerCloudEnvironment.Aws.Required)
			if requiredFieldsNotFilled {
				err = oopsBuilder.
					With("packerCloudEnvironment.Aws.Required", packerCloudEnvironment.Aws.Required).
					Errorf("Required field '%s' is not filled", missingField)
				return
			}
			packerEnvironment.Cloud = packerCloudEnvironment.Aws

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

type CloudEnvironmentI interface {
	IsCloudEnvironment() (isCloudEnvironment bool)
}

type PackerEnvironment struct {
	General *general.PackerGeneralEnvironment
	Cloud   CloudEnvironmentI
}

type Option func(*PackerEnvironment) (err error)

type PackerCloudEnvironmentOptions struct {
	Aws *aws.PackerAwsEnvironment
}