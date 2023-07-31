package cmd

import (
	"github.com/spf13/cobra"
)

type CobraCommands = []*cobra.Command

func GetAllCommands() *CobraCommands {
	return &CobraCommands{
		BuildCommand(),
		UpCommand(),
		DestroyCommand(),
	}
}
