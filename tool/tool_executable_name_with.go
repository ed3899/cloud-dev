package tool

func ToolExecutableNameWith(
	fmtSprintf func(format string, a ...any) string,
) ToolExecutableName {
	toolExecutableName := func(toolName string) string {
		return fmtSprintf(
			"%s.exe",
			toolName,
		)
	}

	return toolExecutableName
}

type ToolExecutableName func(string) string
