package utils

func ToolExecutablePath(
	args *ToolExecutablePathArgs,
) string {
	return args.filepath_Join(
		args.CurrentExecutablePath(),
		args.ToolName(),
		args.ToolExecutableName(),
	)
}

type ToolExecutablePathArgs struct {
	filepath_Join         func(...string) string
	CurrentExecutablePath func() string
	ToolName              func() string
	ToolExecutableName    func() string
}
