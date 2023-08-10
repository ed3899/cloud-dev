package host

import (
	"runtime"

	"github.com/ed3899/kumo/common/functions"
)

func currentArch() functions.Arch {
	return functions.Arch(runtime.GOARCH)
}
