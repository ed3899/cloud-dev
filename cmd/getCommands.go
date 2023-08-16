package cmd

import (
	"github.com/spf13/cobra"
)

type CobraCommands = []*cobra.Command

func GetCommands() CobraCommands {
	return CobraCommands{
		Build(),
		Up(),
		DestroyCommand(),
	}
}
