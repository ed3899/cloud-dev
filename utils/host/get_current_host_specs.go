package host

import (
	"runtime"
)

type GetCurrentHostSpecsF func() (os, arch string)

func GetCurrentHostSpecs() (os, arch string) {
	return runtime.GOOS, runtime.GOARCH
}
