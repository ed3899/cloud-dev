package templates

import (
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/templates"
	"github.com/ed3899/kumo/common/tool"
	packer_aws "github.com/ed3899/kumo/templates/packer/aws"
	packer_general "github.com/ed3899/kumo/templates/packer/general"
	terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	terraform_general "github.com/ed3899/kumo/templates/terraform/general"
	"github.com/samber/oops"
)

func PickTemplate(toolSetup tool.ToolSetupI, cloudSetup cloud.CloudSetupI) (pickedTemplate *MergedTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("pick_template_failed")
		toolType      = toolSetup.GetToolType()
		cloudType     = cloudSetup.GetCloudType()
		packerName    = tool.PACKER_NAME
		terraformName = tool.TERRAFORM_NAME
		awsName       = cloud.AWS_NAME

		generalTemplate, cloudTemplate templates.TemplateSingle
		packerManifest                 templates.PackerManifestI
	)

	// 1. Pick general template
	switch toolType {
	case tool.Packer:
		// 2. Pick general template
		if generalTemplate, err = packer_general.NewTemplate(); err != nil {
			err = oopsBuilder.
				With("tool", tool.Packer).
				Wrapf(err, "Error occurred while picking general template for %s", packerName)
			return
		}
		// 3. Pick cloud template
		switch cloudType {
		case cloud.AWS:
			if cloudTemplate, err = packer_aws.NewTemplate(); err != nil {
				err = oopsBuilder.
					With("tool", tool.Packer).
					With("cloud", cloud.AWS).
					Wrapf(err, "Error occurred while picking template for tool %s and cloud %s", packerName, awsName)
				return
			}

		default:
			err = oopsBuilder.
				With("tool", tool.Packer).
				With("cloud", cloudSetup).
				Wrapf(err, "Error occurred while picking template for tool %s and cloud %v. Unsupported cloud", packerName, cloudSetup)
			return
		}

	case tool.Terraform:
		// 2. Pick general template
		if generalTemplate, err = terraform_general.NewTemplate(); err != nil {
			err = oopsBuilder.
				With("tool", tool.Terraform).
				Wrapf(err, "Error occurred while picking general template for %s", terraformName)
			return
		}
		// 3. Pick cloud template
		switch cloudType {
		case cloud.AWS:
			if packerManifest, err = packer_aws.NewManifest(); err != nil {
				err = oopsBuilder.
					With("tool", tool.Terraform).
					With("cloud", cloud.AWS).
					Wrapf(err, "Error occurred while picking packer manifest for cloud %s", awsName)
				return
			}

			if cloudTemplate, err = terraform_aws.NewTemplate(packerManifest); err != nil {
				err = oopsBuilder.
					With("tool", tool.Terraform).
					With("cloud", cloud.AWS).
					Wrapf(err, "Error occurred while picking template for tool %s and cloud %s", terraformName, awsName)
				return
			}

		default:
			err = oopsBuilder.
				With("tool", tool.Terraform).
				With("cloud", cloudSetup).
				Wrapf(err, "Error occurred while picking template for tool %s and cloud %v. Unsupported cloud", terraformName, cloudSetup)
			return
		}

	default:
		err = oopsBuilder.
			With("tool", toolSetup).
			Wrapf(err, "Error occurred while picking template for tool %v. Unsupported tool", toolSetup)
		return
	}

	// 4. Merge them
	if pickedTemplate, err = NewMergedTemplate(generalTemplate, cloudTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging general template '%v' and cloud template '%v'", generalTemplate.GetAbsPath(), cloudTemplate.GetAbsPath())
		return
	}

	return
}
