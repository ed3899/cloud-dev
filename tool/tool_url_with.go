package tool

func ToolUrlWith[
	ToolName ~func() string,
	ToolVersion ~func() string,
	CurrentOs ~func() string,
	CurrentArch ~func() string,
	Formatter ~func(string, ...interface{}) string,
](
	toolName ToolName,
	toolVersion ToolVersion,
	currentOs CurrentOs,
	currentArch CurrentArch,
	fmt_Sprintf Formatter,
) ToolUrl {
	return func() string {
		return fmt_Sprintf(
			"https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip",
			toolName(),
			toolVersion(),
			toolName(),
			toolVersion(),
			currentOs(),
			currentArch(),
		)
	}
}

type ToolUrl func() string
