package packer

import (
	"github.com/ed3899/kumo/common/cloud"
	"github.com/ed3899/kumo/common/tool"
	"github.com/samber/oops"
)

type Executable struct {
}

func NewExecutable(tool tool.ConfigI, cloud cloud.ConfigI) (executable *Executable, err error) {
	var (
		oopsBuilder = oops.
			Code("new_executable_failed").
			With("tool", tool.GetName()).
			With("version", tool.GetVersion()).
			With("cloud", cloud.GetName())
	)

}
