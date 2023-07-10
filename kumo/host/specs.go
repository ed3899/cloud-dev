package host

import (
	"runtime"
)

type Specs struct {
	OS   string
	ARCH string
}

func GetSpecs() *Specs {
	return &Specs{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}

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
