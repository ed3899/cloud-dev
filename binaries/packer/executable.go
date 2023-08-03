package packer

import (
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
)

type Executable struct {
}

func NewExecutable(tool tool.ConfigI) (executable *Executable, err error) {
	var (
		oopsBuilder = oops.
			Code("new_executable_failed").
			With("tool.GetToolName()", tool.GetToolName()).
			With("tool.GetToolVersion()", tool.GetToolVersion())
	)

}
