package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

func init() {
	if utils.HostIsNotCompatible() {
		var (
			oopsBuilder = oops.
					Code("host_is_not_compatible")
			os, arch = utils.GetCurrentHostSpecs()
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
