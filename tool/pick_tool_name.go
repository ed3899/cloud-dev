package tool

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"

	"github.com/samber/mo"
)

func PickToolName(
	toolIota iota.Tool,
) mo.Result[string] {
	oopsBuilder := oops.
		Code("ToolNameWith").
		With("toolIota", toolIota)

	switch toolIota {
	case iota.Packer:
		return mo.Ok("packer")

	case iota.Terraform:
		return mo.Ok("terraform")

	default:
		err := oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			toolIota,
		)
		return mo.Err[string](err)
	}
}
