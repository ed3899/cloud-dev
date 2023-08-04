package templates

import (
	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	common_cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	"github.com/ed3899/kumo/common/templates"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	common_tool_interfaces "github.com/ed3899/kumo/common/tool/interfaces"
	packer_aws "github.com/ed3899/kumo/templates/packer/aws"
	packer_general "github.com/ed3899/kumo/templates/packer/general"
	terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	terraform_general "github.com/ed3899/kumo/templates/terraform/general"
	"github.com/samber/oops"
)

func New(tool common_tool_interfaces.Tool, cloud common_cloud_interfaces.Cloud) (pickedTemplate *MergedTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("pick_template_failed")
		toolType      = tool.Kind()
		cloudType     = cloud.Kind()
		packerName    = common_tool_constants.PACKER_NAME
		terraformName = common_tool_constants.TERRAFORM_NAME
		awsName       = common_cloud_constants.AWS_NAME

		generalTemplate, cloudTemplate templates.TemplateSingle
		packerManifest                 templates.PackerManifestI
	)

	// 1. Pick general template
	switch toolType {
	case common_tool_constants.Packer:
		// 2. Pick general template
		if generalTemplate, err = packer_general.New(); err != nil {
			err = oopsBuilder.
				With("tool", common_tool_constants.Packer).
				Wrapf(err, "Error occurred while picking general template for %s", packerName)
			return
		}
		// 3. Pick cloud template
		switch cloudType {
		case common_cloud_constants.AWS:
			if cloudTemplate, err = packer_aws.New(); err != nil {
				err = oopsBuilder.
					With("tool", common_tool_constants.Packer).
					With("cloud", common_cloud_constants.AWS).
					Wrapf(err, "Error occurred while picking template for tool %s and cloud %s", packerName, awsName)
				return
			}

		default:
			err = oopsBuilder.
				With("tool", common_tool_constants.Packer).
				With("cloud", cloud).
				Wrapf(err, "Error occurred while picking template for tool %s and cloud %v. Unsupported cloud", packerName, cloud)
			return
		}

	case common_tool_constants.Terraform:
		// 2. Pick general template
		if generalTemplate, err = terraform_general.NewTemplate(); err != nil {
			err = oopsBuilder.
				With("tool", common_tool_constants.Terraform).
				Wrapf(err, "Error occurred while picking general template for %s", terraformName)
			return
		}
		// 3. Pick cloud template
		switch cloudType {
		case common_cloud_constants.AWS:
			if packerManifest, err = packer_aws.NewManifest(); err != nil {
				err = oopsBuilder.
					With("tool", common_tool_constants.Terraform).
					With("cloud", common_cloud_constants.AWS).
					Wrapf(err, "Error occurred while picking packer manifest for cloud %s", awsName)
				return
			}

			if cloudTemplate, err = terraform_aws.NewTemplate(packerManifest); err != nil {
				err = oopsBuilder.
					With("tool", common_tool_constants.Terraform).
					With("cloud", common_cloud_constants.AWS).
					Wrapf(err, "Error occurred while picking template for tool %s and cloud %s", terraformName, awsName)
				return
			}

		default:
			err = oopsBuilder.
				With("tool", common_tool_constants.Terraform).
				With("cloud", cloud).
				Wrapf(err, "Error occurred while picking template for tool %s and cloud %v. Unsupported cloud", terraformName, cloud)
			return
		}

	default:
		err = oopsBuilder.
			With("tool", tool).
			Wrapf(err, "Error occurred while picking template for tool %v. Unsupported tool", tool)
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
