package tool

import (
	"github.com/samber/mo"
	"github.com/samber/oops"
)

func ToolTemplatePathWith(
	filepathJoin func(...string) string,
	fmtSprintf func(string, ...any) string,
	osExecutable func() (string, error),
) mo.Result[ToolTemplatePath] {
	oopsBuilder := oops.
		Code("ToolTemplatePathWith")

	osExecutablePath, err := osExecutable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Failed to get current executable path")
		return mo.Err[ToolTemplatePath](err)
	}

	toolTemplatePath := func(toolName, cloudName, templateName string) string {
		return filepathJoin(
			osExecutablePath,
			toolName,
			cloudName,
			templateName,
		)
	}

	return mo.Ok[ToolTemplatePath](toolTemplatePath)
}

type ToolTemplatePath func(cloudName, toolName, templateName string) string
