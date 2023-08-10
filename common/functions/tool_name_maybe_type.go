package functions

import "github.com/ed3899/kumo/common/iota"

type ToolNameMaybe func(*ToolNameMaybeArgs) ToolName

type ToolName func() string

type ToolNameMaybeArgs struct {
	Tool          iota.Tool
	PackerName    func() string
	TerraformName func() string
}
