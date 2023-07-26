package cmd

import (
	"log"

	"github.com/ed3899/kumo/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	// Read the config
	if err := config.ReadKumoConfig(&config.KumoConfig{
		Name: "kumo.config",
		Type: "yaml",
		Path: ".",
	}); err != nil {
		log.Fatal(errors.Wrapf(err, "Error occurred while reading kumo config"))
	}

	// Assemble commands
	rootCmd.AddCommand(*GetAllCommands()...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(errors.Wrapf(err, "Error occurred while running kumo"))
	}
}
