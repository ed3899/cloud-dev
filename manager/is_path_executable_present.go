package manager

func IsPathExecutablePresentWith(
	utilsFileIsFilePresent func(string) bool,
) IsPathExecutablePresent {
	isPathExecutablePresent := func(manager IPathGetter[IExecutableGetter]) bool {
		return utilsFileIsFilePresent(manager.Path().Executable())
	}

	return isPathExecutablePresent
}

type IsPathExecutablePresent func(IPathGetter[IExecutableGetter]) bool
