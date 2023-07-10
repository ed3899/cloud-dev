package cmd

import (
	"log"

	"github.com/ed3899/kumo/binz"
	"github.com/ed3899/kumo/binz/download"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CobraCommands = []*cobra.Command

func GetAllCommands(bins *download.Binaries) *CobraCommands {
	packer, err := binz.GetPackerInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer instance")
		log.Fatal(err)
	}

	pulumi, err := binz.GetPulumiInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Pulumi instance")
		log.Fatal(err)
	}

	ccmds := []*cobra.Command{
		GetBuildCommand(packer),
		GetUpCommand(pulumi),
		GetDestroyCommand(pulumi),
	}

	return &ccmds
}
