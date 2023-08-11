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

	toolTemplateName := func(cloudName string) string {
		return fmtSprintf("%s.tmpl", cloudName)
	}

	toolTemplatePath := func(cloudName, toolName string) string {
		return filepathJoin(
			osExecutablePath,
			toolName,
			cloudName,
			toolTemplateName(cloudName),
		)
	}
}

type ToolTemplatePath func(cloudName, toolName string) string
