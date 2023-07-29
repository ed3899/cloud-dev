package templates

import (
	"github.com/ed3899/kumo/templates/packer/general"
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
)

func PickTemplate(toolType tool.Type, cloudType cloud.Type) (pickedTemplate *MergedTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("create_template_failed")

		generalTemplate, cloudTemplate TemplateSingle
	)

	// 1. Pick general template
	switch toolType {
	case tool.Packer:
		generalTemplate, err = general.NewTemplate()
		// 2. Pick cloud template
		switch cloudType {
		case cloud.AWS:
		default:
		}
	case tool.Terraform:
		// 2. Pick cloud template
		switch cloudType {
		case cloud.AWS:
		default:
		}
	default:
	}

	// 3. Merge them
	if pickedTemplate, err = NewMergedTemplate(generalTemplate, cloudTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging general template '%v' and cloud template '%v'", generalTemplate.GetName(), cloudTemplate.GetName())
		return
	}

	return
}
