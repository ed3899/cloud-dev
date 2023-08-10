package tool

func ToolExecutablePathWith[
	FilepathJoin ~func(...string) string,
	CurrentExecutablePath ~func() string,
	ToolName ~func() string,
	ToolExecutableName ~func() string,
](
	filepathJoin FilepathJoin,
	currentExecutablePath CurrentExecutablePath,
	toolName ToolName,
	toolExecutableName ToolExecutableName,
) ToolExecutablePath {
	return func() string {
		return filepathJoin(
			currentExecutablePath(),
			toolName(),
			toolExecutableName(),
		)
	}
}

type ToolExecutablePath func() string
