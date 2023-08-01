package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

// TODO Set ssh key file permissions to 600
// TODO add AMI.User tag to packer ami
// TODO add packer clean up
// TODO add ssh config file
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
