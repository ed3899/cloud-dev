package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils"
)

// TODO Set ssh key file permissions to 600
// TODO add AMI.User tag to packer ami
// TODO add cleanup functions for both terraform and packer
func init() {
	if utils.HostIsNotCompatible() {
		log.Fatal("Host is not compatible with kumo :/")
	}
}

func main() {
	cmd.Execute()
}
