package tool

import (
	"fmt"

	"github.com/ed3899/kumo/common/functions"
)

func ToolUrlWith(
	args *functions.ToolUrlWithArgs[func() string],
) functions.ToolName {
	return func() string {
		return fmt.Sprintf(
			"https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip",
			args.ToolName(),
			args.ToolVersion(),
			args.ToolName(),
			args.ToolVersion(),
			args.CurrentOs(),
			args.CurrentArch(),
		)
	}
}
