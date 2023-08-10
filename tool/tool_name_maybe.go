package tool

import (
	"github.com/ed3899/kumo/common/functions"
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ToolNameMaybe(
	args *functions.ToolNameMaybeArgs,
) (
	ToolName func() string,
	err error,
) {
	oopsBuilder := oops.
		Code("PickToolName").
		With("args", *args)

	switch args.Tool {
	case iota.Packer:
		return args.PackerName, nil
	case iota.Terraform:
		return args.TerraformName, nil
	default:
		return nil, oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			args.Tool,
		)
	}
}
