package host

import (
	"runtime"

	"github.com/ed3899/kumo/common/functions"
)

func currentOs() functions.Os {
	return functions.Os(runtime.GOOS)
}
