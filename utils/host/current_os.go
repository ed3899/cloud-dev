package host

import (
	"runtime"
)

func currentOs() string {
	return runtime.GOOS
}