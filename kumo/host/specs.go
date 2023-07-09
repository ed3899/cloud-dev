package host

import "runtime"

type SpecsI interface {
	Valid()
	NotValid()
	Get()
}

type Specs struct {
	OS   string
	ARCH string
}

// Validates host compatibility based on wether packer and pulumi
// are supported by the host OS and ARCH
func (s *Specs) Valid() (valid bool) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386":
		case "amd64":
		default:
			return false
		}
	default:
		return false
	}

	return false
}

// Validates host compatibility based on wether packer and pulumi
// are supported by the host OS and ARCH
func (s *Specs) NotValid() (valid bool) {
	return !s.Valid()
}

// Return the current host specs
func (s *Specs) Get() (specs *Specs) {
	return &Specs{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
}
