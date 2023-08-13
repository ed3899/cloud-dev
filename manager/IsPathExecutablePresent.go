package manager

func IsPathExecutablePresent(
	utilsFileIsFilePresent func(string) bool,
	manager IPathGetter,
) bool {
	return utilsFileIsFilePresent(manager.Path().Executable())
}
