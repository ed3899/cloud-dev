package main

import (
	"log"
	"runtime"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils/host"
	"github.com/samber/oops"
)

func init() {
	if !host.HostIsCompatible(runtime.GOOS, runtime.GOARCH) {
		oopsBuilder := oops.
			Code("HostIsCompatible").
			With("os", runtime.GOOS).
			With("arch", runtime.GOARCH)

		log.Fatalf(
			"%+v",
			oopsBuilder.
				Errorf("Host is not compatible with kumo :/"),
		)
	}
}

func main() {
	cmd.Execute()
}
