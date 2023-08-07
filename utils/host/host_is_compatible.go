package host

import "runtime"

func HostIsCompatible() (compatible bool) {
	switch runtime.GOOS {
	case "windows":

		switch runtime.GOARCH {
		case "386":
			compatible = true

		case "amd64":
			compatible = true

		default:
			compatible = false
		}

	case "darwin":

		switch runtime.GOARCH {
		case "amd64":
			compatible = true

		case "arm64":
			compatible = true

		default:
			compatible = false
		}

	default:
		compatible = false
	}
	return
}
