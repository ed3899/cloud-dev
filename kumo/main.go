package main

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/cmd"
	"github.com/ed3899/kumo/utils"
)

// TODO Remove zips after extraction

var Packer *binz.Packer
var Pulumi *binz.Pulumi

func init() {
	bins, err := utils.GetBinaries()
	if err != nil {
		log.Fatalf("Error occurred while getting binaries: %v", err)
	}
	Packer = binz.GetPacker(bins)
	Pulumi = binz.GetPulumi(bins)
}

func main() {
	cmd.Execute()
}
