package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils/host"
	"github.com/samber/oops"
)

func init() {
	if host.HostIsNotCompatible() {
		var (
			oopsBuilder = oops.
					Code("host_is_not_compatible")
			os, arch = host.GetCurrentHostSpecs()
			err      error
		)

		log.Fatalf(
			"%+v",
			oopsBuilder.
				With("os", os).
				With("arch", arch).
				Wrapf(err, "Host is not compatible with kumo :/"),
		)
	}
}

func main() {
	cmd.Execute()
}
