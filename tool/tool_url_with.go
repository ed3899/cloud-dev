package tool

func ToolUrlWith(
	fmtSprintf func(string, ...any) string,
) ToolUrl {
	toolUrl := func(toolName, toolVersion, currentOs, currentArch string) string {
		return fmtSprintf(
			"https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip",
			toolName,
			toolVersion,
			toolName,
			toolVersion,
			currentOs,
			currentArch,
		)
	}

	return toolUrl
}

type ToolUrl func(toolName, toolVersion, currentOs, currentArch string) string
