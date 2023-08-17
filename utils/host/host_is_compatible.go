package host

import "runtime"

func HostIsCompatible() bool {
	switch runtime.GOOS {
	case "windows":

		switch runtime.GOARCH {
		case "386":
			return true

		case "amd64":
			return true

		default:
			return false
		}

	case "darwin":

		switch runtime.GOARCH {
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
