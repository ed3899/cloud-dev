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

type Option func(Tool) (Tool, error)

func WithKind(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("toolKind", toolKind)
	)

	option = func(t Tool) (tool Tool, err error) {
		switch toolKind {
		case constants.Packer:
			t.Kind = constants.Packer
		case constants.Terraform:
			t.Kind = constants.Terraform
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithName(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("toolKind", toolKind)
	)

	option = func(t Tool) (tool Tool, err error) {
		switch toolKind {
		case constants.Packer:
			t.Name = constants.PACKER
		case constants.Terraform:
			t.Name = constants.TERRAFORM
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithVersion(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithVersion").
			With("toolKind", toolKind)
	)

	option = func(t Tool) (tool Tool, err error) {
		switch toolKind {
		case constants.Packer:
			t.Version = constants.PACKER_VERSION
		case constants.Terraform:
			t.Version = constants.TERRAFORM_VERSION
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithUrl(
	toolKind constants.ToolKind,
	createHashicorpUrl utils.CreateHashicorpURLF,
	getCurrentHostSpecs utils.GetCurrentHostSpecsF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithUrl").
				With("toolKind", toolKind)
		currentOs, currentArch = getCurrentHostSpecs()
	)

	option = func(t Tool) (tool Tool, err error) {
		switch toolKind {
		case constants.Packer:
			t.Url = createHashicorpUrl(constants.PACKER, constants.PACKER_VERSION, currentOs, currentArch)
		case constants.Terraform:
			t.Url = createHashicorpUrl(constants.TERRAFORM, constants.TERRAFORM_VERSION, currentOs, currentArch)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
		}

		return
	}

	return
}

func NewTool(opts ...Option) (tool Tool, err error) {
	var (
		oopsBuilder = oops.
				Code("NewTool").
				With("opts", opts)

		o Option
	)

	tool = Tool{}
	for _, o = range opts {
		if tool, err = o(tool); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", o)
			return
		}
	}

	return
}
