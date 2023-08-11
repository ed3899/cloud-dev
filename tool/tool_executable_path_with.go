package tool

import (
	"github.com/samber/mo"
	"github.com/samber/oops"
)

func ToolExecutablePathWith(
	filepathJoin func(...string) string,
	osExecutable func() (string, error),
) mo.Result[ToolExecutablePath] {
	oopsBuilder := oops.
		Code("ToolExecutablePathWith")

	executablePath, err := osExecutable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Failed to get current executable path")
		return mo.Err[ToolExecutablePath](err)
	}

	toolExecutablePath := func(toolName string, toolExecutableName string) string {
		return filepathJoin(
			executablePath,
			toolName,
			toolExecutableName,
		)
	}

	return mo.Ok[ToolExecutablePath](toolExecutablePath)
}

type ToolExecutablePath func(string, string) string
