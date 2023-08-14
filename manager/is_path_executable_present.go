package manager

func IsPathExecutablePresentWith(
	utilsFileIsFilePresent func(string) bool,
) IsPathExecutablePresent {
	isPathExecutablePresent := func(manager IPathGetter) bool {
		return utilsFileIsFilePresent(manager.Path().Executable())
	}

	return isPathExecutablePresent
}

type IsPathExecutablePresent func(IPathGetter) bool
