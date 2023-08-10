package host

import (
	"runtime"
)

func currentArch() string {
	return runtime.GOARCH
}
