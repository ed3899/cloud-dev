package cmd

import (
	"log"

	"github.com/ed3899/kumo/config"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	var (
		oopsBuilder = oops.Code("root_cmd_init_failed")
	)

	// Read the config
	if err := config.ReadKumoConfig(&config.KumoConfig{
		Name: "kumo.config",
		Type: "yaml",
		Path: ".",
	}); err != nil {
		log.Fatalf(
			"%+v",
			oopsBuilder.
				Wrapf(err, "Error occurred while reading kumo config"),
		)
	}

	// Assemble commands
	rootCmd.AddCommand(*GetAllCommands()...)
}

func Execute() {
	var (
		oopsBuilder = oops.Code("root_cmd_execute_failed")
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(
			"%+v",
			oopsBuilder.
				Wrapf(err, "Error occurred while running kumo"),
		)
	}
}
