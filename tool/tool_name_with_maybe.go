package tool

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func ToolNameWithMaybe[
	ToolName ~func() string,
](
	toolIota iota.Tool,
	PackerName ToolName,
	TerraformName ToolName,
) (
	ToolName,
	error,
) {

	oopsBuilder := oops.
		Code("ToolNameWithMaybe").
		With("toolIota", toolIota)

	switch toolIota {
	case iota.Packer:
		return PackerName, nil
	case iota.Terraform:
		return TerraformName, nil
	default:
		return nil, oopsBuilder.Errorf(
			"Unknown tool '%#v'",
			toolIota,
		)
	}
}
