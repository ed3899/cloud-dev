package host

import (
	"runtime"
)

func IsCompatible() bool {
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
	default:
		return false
	}
}

func IsNotCompatible() bool {
	return !IsCompatible()
}
