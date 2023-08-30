package host

// Returns true if the host is compatible with the given OS and architecture.
func HostIsCompatible(os, arch string) bool {
	switch os {
	case "windows":

		switch arch {
		case "386":
			return true

		case "amd64":
			return true

		default:
			return false
		}

	case "darwin":

		switch arch {
		case "amd64":
			return true

		case "arm64":
			return true

		default:
			return false
		}

	default:
		return false
	}
}
