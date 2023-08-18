package cmd

import (
	"github.com/spf13/cobra"
)

type CobraCommands = []*cobra.Command

func Kumo() *cobra.Command {
	kumo := &cobra.Command{
		Use:   "kumo",
		Short: "🌩️ Your quick and easy cloud development environment.",
		Long:  `🌩️ Your quick and easy cloud development environment.`,
	}

	kumo.AddCommand(Subcommands()...)

	return kumo
}

func Subcommands() CobraCommands {
	return CobraCommands{
		Build(),
		Up(),
		Destroy(),
	}
}
