package functions

import "github.com/ed3899/kumo/common/iota"

type ToolNameMaybeArgs struct {
	Tool          iota.Tool
	PackerName    func() string
	TerraformName func() string
}
