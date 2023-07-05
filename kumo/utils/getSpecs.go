package utils

import (
	"runtime"
)

type IHost interface {
	GetSpecs()
}

type Host struct {
}

type Specs struct {
	OS string
	ARCH string
}

func (h *Host) GetSpecs() Specs {
	return Specs{
		OS: runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}

