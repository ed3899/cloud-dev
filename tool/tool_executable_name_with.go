package tool

func ToolExecutableNameWith[
	Fmt_Sprintf ~func(format string, a ...any) string,
	ToolName ~func() string,
](
	fmt_Sprintf Fmt_Sprintf,
	toolName ToolName,
) ToolExecutableName {
	return func() string {
		return fmt_Sprintf(
			"%s.exe",
			toolName(),
		)
	}
}

type ToolExecutableName func() string
