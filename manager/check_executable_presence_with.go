package manager

func CheckExecutablePresenceWith(
	utilsIsFilePresent func(string) bool,
) ForPathGetter {
	return func(m IPathGetter) bool {
		return utilsIsFilePresent(m.Path().Executable())
	}
}

type ForPathGetter func(IPathGetter) bool
