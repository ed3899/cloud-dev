package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
)

// TODO Set ssh key file permissions to 600
// TODO add AMI.User tag to packer ami
// TODO add cleanup functions for both terraform and packer
func init() {
	if utils.HostIsNotCompatible() {
		var (
			os, arch = utils.GetCurrentHostSpecs()
			err      error
		)

		err = oops.Code("host_is_not_compatible").
			With("os", os).
			With("arch", arch).
			Errorf("Host is not compatible with kumo :/")
		log.Fatalf("%+v", err)
	}
}

func main() {
	cmd.Execute()
}
