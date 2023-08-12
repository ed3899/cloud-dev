package manager

func IsExecutablePresentWith(
	utilsFileIsFilePresent func(string) bool,
) ForPathGetter {
	return func(m IPathGetter) bool {
		return utilsFileIsFilePresent(m.Path().Executable())
	}
}

type ForPathGetter func(IPathGetter) bool
