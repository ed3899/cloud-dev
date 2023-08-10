package tool

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"

	"github.com/samber/mo"
)

func ToolNameWith[
	ToolName ~func() string,
](
	toolIota iota.Tool,
	packerName ToolName,
	terraformName ToolName,
) mo.Result[ToolName] {
	oopsBuilder := oops.
		Code("ToolNameWith").
		With("toolIota", toolIota)

	switch toolIota {
	case iota.Packer:
		return mo.Ok[ToolName](packerName)

	case iota.Terraform:
		return mo.Ok[ToolName](packerName)

	default:
		err := oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			toolIota,
		)
		return mo.Err[ToolName](err)
	}
}
