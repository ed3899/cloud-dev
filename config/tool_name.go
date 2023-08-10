package config

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ToolName(
	args *ToolNameArgs,
) (
	string,
	error,
) {
	oopsBuilder := oops.
		Code("PickToolName").
		With("args", *args)

	switch args.Tool {
	case iota.Packer:
		return args.PackerName(), nil
	case iota.Terraform:
		return args.TerraformName(), nil
	default:
		return "", oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			args.Tool,
		)
	}
}

type ToolNameArgs struct {
	Tool          iota.Tool
	PackerName    func() string
	TerraformName func() string
}
