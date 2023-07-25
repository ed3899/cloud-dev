package cmd

import (
	"github.com/spf13/cobra"
)

type CobraCommands = []*cobra.Command

func GetAllCommands() *CobraCommands {
	ccmds := []*cobra.Command{
		GetBuildCommand(),
		GetUpCommand(),
		GetDestroyCommand(),
	}

	return &ccmds
}
