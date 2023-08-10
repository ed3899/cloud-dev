package utils

func ToolExecutableName(
	fmt_Sprintf func(format string, a ...any) string,
	ToolName func() string,
) string {
	return fmt_Sprintf("%s.exe", ToolName())
}
