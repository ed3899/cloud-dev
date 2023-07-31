package hashicorp_vars

import (
	"os"

	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/tool"
	terraform_aws "github.com/ed3899/kumo/hashicorp_vars/terraform/aws"
	packer_aws "github.com/ed3899/kumo/hashicorp_vars/packer/aws"
	"github.com/samber/oops"
)

type HashicorpVarsI interface {
	GetFile() *os.File
}

func PickHashicorpVars(toolType tool.ToolType, cloudType cloud.CloudType) (pickedHashicorpVars HashicorpVarsI , err error) {
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
