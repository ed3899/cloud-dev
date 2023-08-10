package host

import (
	"runtime"
)

func currentOs() string {
	return runtime.GOOS
}

func currentArch() string {
	return runtime.GOARCH
}