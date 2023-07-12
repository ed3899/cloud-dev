package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/host"
)

// TODO Remove zips after extraction
// TODO Set ssh key file permissions to 600
// TODO nest terraform and packer config files in a directory according to the cloud they are deploying to
// TODO add AMI.User tag to packer ami
func init() {
	if host.IsNotCompatible() {
		log.Fatal("Host is not compatible with kumo :/")
	}
}

func main() {
	cmd.Execute()
}
