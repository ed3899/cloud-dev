package utils

import (
	"runtime"
)

func GetCurrentHostSpecs() (os, arch string) {
	return runtime.GOOS, runtime.GOARCH
}

// TODO add darwin
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

func HostIsNotCompatible() (notCompatible bool) {
	return !HostIsCompatible()
}
