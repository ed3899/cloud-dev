package tool

import (
	constants "github.com/ed3899/kumo/0_constants"
	utils "github.com/ed3899/kumo/1_utils"
	"github.com/samber/oops"
)

type Tool struct {
	Name    string
	Version string
	Url     string
}

func NewTool(kind constants.ToolKind) (tool Tool, err error) {
	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed").
				With("tool", kind)
		currentOs, currentArch = utils.GetCurrentHostSpecs()
	)

	switch kind {
	case constants.Packer:
		tool = Tool{
			Name:    constants.PACKER,
			Version: constants.PACKER_VERSION,
			Url: utils.CreateHashicorpURL(constants.PACKER, constants.PACKER_VERSION, currentOs, currentArch),
		}

	case constants.Terraform:
		tool = Tool{
			Name:    constants.TERRAFORM,
			Version: constants.TERRAFORM_VERSION,
			Url: utils.CreateHashicorpURL(constants.TERRAFORM, constants.TERRAFORM_VERSION, currentOs, currentArch),
		}

	default:
		err = oopsBuilder.
			Errorf("Unknown tool kind: %d", kind)
		return

	}

	return
}
