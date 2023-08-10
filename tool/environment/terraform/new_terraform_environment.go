package environment

import (
	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool_manager/environment/terraform/general"
	"github.com/ed3899/kumo/utils/environment"
	"github.com/samber/oops"
)

func NewTerraformEnvironment(
	options ...Option,
) (terraformEnvironment *TerraformEnvironment, err error) {
	var (
		oopsBuilder = oops.
				Code("NewTerraformEnvironment").
				With("options", options)

		opt Option
	)

	terraformEnvironment = &TerraformEnvironment{}
	for _, opt = range options {
		if err = opt(terraformEnvironment); err != nil {
			err = oopsBuilder.
				With("terraformEnvironment", terraformEnvironment).
				Wrapf(err, "Error while applying option %v", opt)
			return
		}
	}

	return
}

func WithGeneralEnvironment(
	terraformGeneralEnvironment *general.TerraformGeneralEnvironment,
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithGeneralEnvironment").
				With("terraformGeneralEnvironment", terraformGeneralEnvironment).
				With("areRequiredFieldsNotFilled", areRequiredFieldsNotFilled)

		requiredFieldsNotFilled bool
		missingField            string
	)

	option = func(terraformEnvironment *TerraformEnvironment) (err error) {
		requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(terraformGeneralEnvironment.Required)
		if requiredFieldsNotFilled {
			err = oopsBuilder.
				With("terraformGeneralEnvironment.Required", terraformGeneralEnvironment.Required).
				Errorf("Required field '%s' is not filled", missingField)
			return
		}

		terraformEnvironment.General = terraformGeneralEnvironment

		return
	}

	return
}

func WithCloudEnvironment(
	cloud cloud.Cloud,
	terraformCloudEnvironmentOptions *TerraformCloudEnvironmentOptions,
	areRequiredFieldsNotFilled environment.IsStructNotCompletelyFilledF,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
				Code("WithCloudEnvironment").
				With("terraformCloudEnvironment", terraformCloudEnvironmentOptions).
				With("areRequiredFieldsNotFilled", areRequiredFieldsNotFilled)

		requiredFieldsNotFilled bool
		missingField            string
	)

	option = func(terraformEnvironment *TerraformEnvironment) (err error) {
		switch cloud.Kind {
		case constants.Aws:
			requiredFieldsNotFilled, missingField = areRequiredFieldsNotFilled(terraformCloudEnvironmentOptions.Aws.Required)
			if requiredFieldsNotFilled {
				err = oopsBuilder.
					With("terraformCloudEnvironmentOptions.Aws.Required", terraformCloudEnvironmentOptions.Aws.Required).
					Errorf("Required field '%s' is not filled", missingField)
				return
			}

			terraformEnvironment.Cloud = terraformCloudEnvironmentOptions.Aws

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

type TerraformEnvironment struct {
	General *general.TerraformGeneralEnvironment
	Cloud   interfaces.TerraformCloudEnvironmentI
}

type Option func(*TerraformEnvironment) (err error)
