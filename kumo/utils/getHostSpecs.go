package utils

import "runtime"

func GetCurrentHostSpecs() (os, arch string) {
	return runtime.GOOS, runtime.GOARCH
}
