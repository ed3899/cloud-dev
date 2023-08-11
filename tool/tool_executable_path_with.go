package tool

import (
	"github.com/samber/mo"
	"github.com/samber/oops"
)

func ToolExecutablePathWith(
	filepathJoin func(...string) string,
	osExecutable func() (string, error),
	fmtSprintf func(string, ...any) string,
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
			fmtSprintf("%s.exe", toolName),
		)
	}

	return mo.Ok[ToolExecutablePath](toolExecutablePath)
}

type ToolExecutablePath func(toolName string, toolExecutableName string) string
