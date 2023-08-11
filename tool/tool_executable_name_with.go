package tool

import (
	"fmt"

	"github.com/ed3899/kumo/common/iota"
)

func ToolExecutableNameWith(
	fmt_Sprintf func(format string, a ...any) string,
) ToolExecutableName {
	return func(toolName string) string {
		return fmt_Sprintf(
			"%s.exe",
			toolName,
		)
	}
}

type ToolExecutableName func(string) string

func T() {


}
