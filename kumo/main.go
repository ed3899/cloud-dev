package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/host"
)

// TODO Set ssh key file permissions to 600
// TODO add AMI.User tag to packer ami
// TODO integration test for packer and terraform urls
func init() {
	if host.IsNotCompatible() {
		log.Fatal("Host is not compatible with kumo :/")
	}
}

func main() {
	cmd.Execute()
}
