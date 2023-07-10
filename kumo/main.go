package main

import (
	"log"

	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/host"
)

// TODO Remove zips after extraction
func init() {
	if host.IsNotCompatible() {
		log.Fatal("Host is not compatible with kumo :/")
	}
}

func main() {
	cmd.Execute()
}
