package manager

func CheckExecutablePresenceWith(
	utilsIsFilePresent func(string) bool,
) ForManager {
	return func(m Manager) bool {
		return utilsIsFilePresent(m.Path().Executable())
	}
}

type ForManager func(Manager) bool
