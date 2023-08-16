package main

import (
	"log"
	"runtime"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils/host"
	"github.com/samber/oops"
)

func init() {
	if host.HostIsNotCompatible() {
		oopsBuilder := oops.
			Code("host_is_not_compatible")

		log.Fatalf(
			"%+v",
			oopsBuilder.
				With("os", runtime.GOOS).
				With("arch", runtime.GOARCH).
				Errorf("Host is not compatible with kumo :/"),
		)
	}
}

func main() {
	cmd.Execute()
}
