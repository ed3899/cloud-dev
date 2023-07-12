package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/host"
)

// TODO Remove zips after extraction
// TODO Set ssh key file permissions to 600
func init() {
	if host.IsNotCompatible() {
		log.Fatal("Host is not compatible with kumo :/")
	}
}

func main() {
	cmd.Execute()
}
