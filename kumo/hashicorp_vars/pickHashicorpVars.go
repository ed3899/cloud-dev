package hashicorp_vars

import (
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/tool"
	packer_aws "github.com/ed3899/kumo/hashicorp_vars/packer/aws"
	terraform_aws "github.com/ed3899/kumo/hashicorp_vars/terraform/aws"
	"github.com/samber/oops"
)

func PickHashicorpVars(toolType tool.ToolType, cloudType cloud.CloudType) (pickedHashicorpVars hashicorp_vars.HashicorpVarsI, err error) {
	var (
		oopsBuilder = oops.
			Code("pick_hashicorp_vars_failed").
			With("toolType", toolType).
			With("cloudType", cloudType)
	)

	switch toolType {
	case tool.Packer:

		switch cloudType {
		case cloud.AWS:
			if pickedHashicorpVars, err = packer_aws.NewHashicorpVars(); err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while instantiating Packer AWS Hashicorp Vars")
				return
			}

		default:
			err = oopsBuilder.
				Errorf("Cloud '%v' not supported", cloudType)
			return
		}

	case tool.Terraform:

		switch cloudType {
		case cloud.AWS:
			if pickedHashicorpVars, err = terraform_aws.NewHashicorpVars(); err != nil {
				err = oopsBuilder.
					Wrapf(err, "Error occurred while instantiating Terraform AWS Hashicorp Vars")
				return
			}

		default:
			err = oopsBuilder.
				Errorf("Cloud '%v' not supported", cloudType)
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolType)
		return
	}

	return
}
