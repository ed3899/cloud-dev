package tool

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ToolNameMaybeWith[
	ToolName ~func() string,
](
	toolIota iota.Tool,
	packerName ToolName,
	terraformName ToolName,
) (
	ToolName,
	error,
) {

	oopsBuilder := oops.
		Code("ToolNameMaybeWith").
		With("toolIota", toolIota)

	switch toolIota {
	case iota.Packer:
		return packerName, nil
	case iota.Terraform:
		return terraformName, nil
	default:
		return nil, oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			toolIota,
		)
	}
}
