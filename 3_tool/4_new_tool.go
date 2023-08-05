package tool

import (
	constants "github.com/ed3899/kumo/constants"
	utils "github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

type Tool struct {
	Kind    constants.ToolKind
	Name    string
	Version string
	Url     string
}

func NewTool(kind constants.ToolKind) (ForSpecs ForSpecs) {
	var (
		oopsBuilder = oops.
			Code("new_tool_setup_failed").
			With("tool", kind)
	)

	ForSpecs = func(currentOs, currentArch string) (tool Tool, err error) {
		switch kind {
		case constants.Packer:
			tool = Tool{
				Kind:    constants.Packer,
				Name:    constants.PACKER,
				Version: constants.PACKER_VERSION,
				Url:     utils.CreateHashicorpURL(constants.PACKER, constants.PACKER_VERSION, currentOs, currentArch),
			}

		case constants.Terraform:
			tool = Tool{
				Kind:    constants.Terraform,
				Name:    constants.TERRAFORM,
				Version: constants.TERRAFORM_VERSION,
				Url:     utils.CreateHashicorpURL(constants.TERRAFORM, constants.TERRAFORM_VERSION, currentOs, currentArch),
			}

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", kind)
			return

		}

		return
	}

	return
}

type ForSpecs = func(currentOs, currentArch string) (tool Tool, err error)
